package cart

import (
	"warehouse-application/pkg/database"
	"warehouse-application/pkg/logging"
)

//структура карточки

type Cart struct {
	ProductID   uint `json:"productID"`
	UserID      uint
	ProductName string `json:"productName"`
	Price       uint   `json:"price"`
}

//функция добавления новой карточки

func (c *Cart) AppendCartToDB(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `INSERT INTO product(user_id, name, price) VALUES($1,$2,$3)`

	if _, err := db.Exec(data, c.UserID, c.ProductName, c.Price); err != nil {
		logger.Info(err)
		return
	}
	logger.Info("cart has been added to DB")
}

//структура для удобного вывода всех карточек пользователя

type SuperCart struct {
	Carts []struct {
		ProductID   uint
		ProductName string
		Price       uint
	}
}

//функция для показа всех карточек которые есть у пользователя

func ShowCartsOfUser(userId uint, logger logging.Logger) SuperCart {
	db := database.ConnectToDB()
	defer db.Close()

	var result SuperCart

	data := `SELECT id,name,price FROM product WHERE user_id = $1`

	query, err := db.Query(data, userId)

	if err != nil {
		logger.Info(err)
		return SuperCart{nil}
	}

	defer query.Close()

	for query.Next() {

		cart := struct {
			ProductID   uint
			ProductName string
			Price       uint
		}{}

		err = query.Scan(&cart.ProductID, &cart.ProductName, &cart.Price)
		if err != nil {
			logger.Info(err)
			return SuperCart{nil}
		}
		result.Carts = append(result.Carts, cart)
	}
	if query.Err() != nil {
		logger.Info(query.Err())
		return SuperCart{nil}
	}

	return result
}

func (c *Cart) CorrectPrice(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `UPDATE product SET price = $1 WHERE id = $2 and user_id = $3`

	if _, err := db.Exec(data, c.Price, c.ProductID, c.UserID); err != nil {
		logger.Info(err)
		return
	}
	logger.Info("price has been corrected")
}

func (c *Cart) DeleteCart(logger logging.Logger) {
	db := database.ConnectToDB()
	defer db.Close()

	data := `DELETE FROM stock WHERE product_id = $1 and user_id = $2`

	if _, err := db.Exec(data, c.ProductID, c.UserID); err != nil {
		logger.Info(err)
		return
	}

	data = `DELETE FROM product WHERE id = $1`

	if _, err := db.Exec(data, c.ProductID); err != nil {
		logger.Info(err)
		return
	}

	logger.Info("cart has been deleted")
}
