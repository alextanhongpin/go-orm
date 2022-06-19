package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/alextanhongpin/orm"
	_ "github.com/lib/pq"
)

/** Migrations

create table users (
	id int generated always as identity,
	name text not null,
	primary key (id),
	unique (name)
);

create table books (
	id int generated always as identity,
	title text,
	user_id int not null,
	primary key (id),
	foreign key (user_id) references users(id)
);

insert into users(name) values ('john appleseed');
insert into books(title, user_id) values ('the meaning of life', 1);
*/

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	dbname   = os.Getenv("DB_NAME")
)

type User struct {
	ID   int    `sql:"id"`
	Name string `sql:"name"`
}

func (u *User) Fields() []any {
	return []any{&u.ID, &u.Name}
}

type Book struct {
	ID     int    `sql:"id"`
	Title  string `sql:"title"`
	UserID int    `sql:"user_id"`
}

func (b *Book) Fields() []any {
	return []any{&b.ID, &b.Title, &b.UserID}
}

var UserDataMapper = orm.NewDataMapper[User]("users", "u", "sql")
var BookDataMapper = orm.NewDataMapper[Book]("books", "b", "sql")

func main() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	id := 1

	u := &User{}
	if err := db.
		QueryRow(`SELECT * FROM users WHERE id = $1`, id).
		Scan(u.Fields()...); err != nil {
		log.Fatalf("failed to query row: %s", err)
	}
	fmt.Printf("%+v\n", u)

	{
		b := &Book{}
		u := &User{}

		// SELECT * is prone to errors, especially when new columns are
		// added/removed.
		fields := append(u.Fields(), b.Fields()...)
		stmt, args := orm.WhereStmt(`
				SELECT u.*, b.*
				FROM users u
				JOIN books b ON (b.user_id = u.id)
			`,
			[]*orm.Wherer{
				orm.Where("u.id", 1),
				orm.Where("b.user_id", 1),
			},
		)

		fmt.Println("stmt:", stmt)
		if err := db.
			QueryRow(stmt, args...).
			Scan(fields...); err != nil {
			log.Fatalf("failed to query row: %s", err)
		}
		fmt.Printf("user: %+v\n", u)
		fmt.Printf("book: %+v\n", b)
	}

	{
		u := &User{}
		stmt := fmt.Sprintf(
			`
				SELECT %s 
				FROM %s 
				WHERE id = $1
			`,
			UserDataMapper.Columns(),
			UserDataMapper.SelectName(),
		)
		args := []any{1}

		fmt.Println("stmt:", stmt)

		if err := db.
			QueryRow(stmt, args...).
			Scan(UserDataMapper.Fields(u)...); err != nil {
			log.Fatalf("failed to query row: %s", err)
		}

		fmt.Printf("user: %+v\n", u)
	}

	{
		u := &User{}
		b := &Book{}

		stmt := fmt.Sprintf(
			`
				SELECT %s, %s 
				FROM users u
				JOIN books b ON (b.user_id = u.id)
			`,
			UserDataMapper.Columns(),
			BookDataMapper.Columns(),
		)
		stmt, args := orm.WhereStmt(stmt, []*orm.Wherer{
			orm.Where("u.id", 1),
		})
		cols := append(UserDataMapper.Fields(u), BookDataMapper.Fields(b)...)

		fmt.Println("stmt:", stmt)
		fmt.Println("args:", args)

		if err := db.
			QueryRow(stmt, args...).
			Scan(cols...); err != nil {
			log.Fatalf("failed to query row: %s", err)
		}

		fmt.Printf("user: %+v\n", u)
		fmt.Printf("book: %+v\n", b)
	}

	{
		// Due to varying args length, it becomes tricky when all you need is to
		// set a function, or use string functions like LOWER(). (hint: just to it
		// at app layer).
		stmt, args := orm.UpdateStmt(
			"UPDATE users",
			orm.Set("name", "john appleseed"),
			//orm.Set("updated_at = now()"),
		)

		stmt, args = orm.WhereStmt(
			stmt,
			[]*orm.Wherer{
				orm.Where("id", 1),
				orm.Where("1 = 1"),
			},
			args...,
		)

		fmt.Println("stmt:", stmt)
		fmt.Println("args:", args)

		if res, err := db.
			Exec(stmt, args...); err != nil {
			log.Fatalf("failed to query row: %s", err)
		} else {
			fmt.Println("did update?")
			fmt.Println(res.RowsAffected())
		}
	}
}
