package entity

type Person struct {
	ID      uint64  `db:"id"`
	Name    string  `db:"name"`
	Age     *int    `db:"age"`
	Address *string `db:"address"`
	Work    *string `db:"work"`
}
