package models

type Product struct {
	Id    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Stock int    `db:"stock" json:"stock"`
}
