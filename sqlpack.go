package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

const testSuffix = "_test"

const (
	flagWatch    = "watch"
	flagFile     = "file"
	flagPostgres = "postgres"
	flagStdout   = "stdout"
)

var version string

func main() {
	cliFlagWatch := cli.BoolFlag{
		Name:  fmt.Sprintf("%s, w", flagWatch),
		Usage: "watch for file changes",
	}
	cliFlagPostgres := cli.StringFlag{
		Name:  fmt.Sprintf("%s, pg", flagPostgres),
		Value: "",
		Usage: "postgres connection string",
	}
	app := cli.NewApp()
	app.Name = "sqlpack"
	app.Version = version
	app.Usage = "bundle files and deploy them to database"
	app.Flags = []cli.Flag{
		cliFlagWatch,
		cli.BoolFlag{
			Name:  fmt.Sprintf("%s, std", flagStdout),
			Usage: "write bundle into console",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s, f", flagFile),
			Value: "",
			Usage: "output file path",
		},
		cliFlagPostgres,
	}
	app.Action = action
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "test",
			Usage:  "execute test files one by one",
			Action: testAction,
			Flags:  []cli.Flag{cliFlagWatch, cliFlagPostgres},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func action(c *cli.Context) (err error) {
	path := c.Args().First()

	type textFunc func(text string) error
	var funcs []textFunc
	if c.IsSet(flagFile) {
		funcs = append(funcs, func(text string) (err error) {
			err = ioutil.WriteFile(c.String(flagFile), []byte(text), os.ModePerm)
			return
		})
	}
	if c.IsSet(flagStdout) {
		funcs = append(funcs, func(text string) (err error) {
			fmt.Println(text)
			return
		})
	}
	if c.IsSet(flagPostgres) {
		funcs = append(funcs, func(text string) (err error) {
			var db *sql.DB
			db, err = sql.Open("postgres", c.String(flagPostgres))
			defer db.Close()
			if err != nil {
				return
			}
			_, err = db.Exec(text)
			return
		})
	}

	exec := func(path string) {
		text, err := bundle(path)
		if err != nil {
			log.Println(err)
		}
		for _, fn := range funcs {
			err := fn(text)
			if err != nil {
				log.Println(err)
			}
		}
	}
	exec(path)
	if c.Bool(flagWatch) {
		err = watch(filepath.Dir(path), exec)
	}
	return
}

func testAction(c *cli.Context) (err error) {
	if !c.IsSet(flagPostgres) {
		log.Fatalln("no db connection provided")
	}
	db, err := openPostgres(c.String(flagPostgres))
	if err != nil {
		log.Fatalln(err)
	}
	defer db.close()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	exec := func(path string) {
		l := newStdLog(path)
		text, err := bundle(path)
		if err != nil {
			l.fail(err.Error())
			return
		}
		db.execLog(text, l)
	}

	dir := c.Args().First()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if isTestFile(path) {
			exec(path)
		}
		return nil
	})
	if err != nil {
		return
	}

	if c.Bool(flagWatch) {
		err = watch(dir, func(path string) {
			if filepath.IsAbs(path) {
				path, err = filepath.Rel(wd, path)
				if err != nil {
					log.Println(err)
					return
				}
			}
			if !isTestFile(path) {
				ext := filepath.Ext(path)
				path = strings.TrimSuffix(path, ext) + testSuffix + ext
			}
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return
			}
			exec(path)
		})
		if err != nil {
			return
		}
	}
	return
}

func isTestFile(path string) bool {
	return strings.HasSuffix(path, testSuffix+filepath.Ext(path))
}
