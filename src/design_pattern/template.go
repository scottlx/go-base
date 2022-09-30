package design_pattern

import "fmt"

//模板方法模式

type Dish interface {
	OpenFire()
	AddOil()
	AddIngredient()
	Stir()
	AddSeasoning()
}

type template struct {
	d Dish
}

func (t template) MakeDish() {
	if t.d == nil {
		return
	}
	fmt.Println("========start to cook=========")
	t.d.OpenFire()
	fmt.Println("file opened")
	t.d.AddOil()
	fmt.Println("oil added, wait until hot")
	t.d.AddIngredient()
	fmt.Println("food added, start stir")
	t.d.Stir()
	fmt.Println("time to add seasoning")
	t.d.AddSeasoning()
}

type Fish struct {
	template
}

func NewFish() *Fish {
	fish := &Fish{}
	fish.d = fish
	return fish
}

func (f Fish) OpenFire() {
	fmt.Println("set fire to big")
}

func (f Fish) AddOil() {
	fmt.Println("add oliver oil")
}

func (f Fish) AddIngredient() {
	fmt.Println("add fish")
}

func (f Fish) Stir() {
	fmt.Println("no stir")
}

func (f Fish) AddSeasoning() {
	fmt.Println("add salt, MSG and soybean sauce")
}

type ChickBreast struct {
	template
}

func NewChickBreast() *ChickBreast {
	br := &ChickBreast{}
	br.d = br
	return br
}

func (c ChickBreast) OpenFire() {
	fmt.Println("set fire to really big")
}

func (c ChickBreast) AddOil() {
	fmt.Println("add oliver oil")
}

func (c ChickBreast) AddIngredient() {
	fmt.Println("add sliced breast")
}

func (c ChickBreast) Stir() {
	fmt.Println("stir heavily")
}

func (c ChickBreast) AddSeasoning() {
	fmt.Println("add salt, MSG and pepper")
}
