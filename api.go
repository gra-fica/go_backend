package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
)

type Claims struct {
	name string
	jwt.RegisteredClaims
}

func bind_apis(e *echo.Echo, database *Database) {
	// non private apis
	e.GET("/api/v1/product/list", func (c echo.Context) error {
		products, err := database.ListProducts();
		if err != nil {
			return c.String(500, err.Error());
		}
		type ProductResponse struct {
			Products []Product `json:"products"`
		}
		return c.JSON(200, ProductResponse{products})
	});

	e.GET("/api/v1/product/search/:name", func (c echo.Context) error {
		name := c.Param("name");

		products, err := database.SearchProduct(name, &ListAllProducts{}, &TokenFuzzy{});
		if err != nil {

			return c.String(500, err.Error());
		}
		type SearchProductResponse struct {
			Products []ProductMatch `json:"products"`
		}
		return c.JSON(200, SearchProductResponse{products});
	});

	e.GET("/api/v1/product/search_exact/:name", func (c echo.Context) error {
		name := c.Param("name");
		product, err := database.GetProductFromName(name);
		if err != nil{
			return c.String(404, fmt.Sprintf("%s not found", name));
		}
		return c.JSON(200, product);
	})

	e.GET("/api/v1/product/search/:name/price/:price", func (c echo.Context) error {
		name := c.Param("name");
		price, err := strconv.Atoi(c.Param("price"));
		if err != nil {
			return c.String(400, fmt.Sprintf("price {%v} is not a number", c.Param("price")));
		}

		products, err := database.SearchProduct(name, &ListProductsWherePrice{price}, &TokenFuzzy{});
		type SearchProductResponse struct {
			Products []ProductMatch `json:"products"`
		}
		return c.JSON(200, SearchProductResponse{products})
	})

	// need auth apis
	// todo convert into QueryParam and require auth
	e.DELETE("/api/v1/product/delete/:id", func (c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64);
		if err != nil {
			return c.String(400, fmt.Sprintf("id {%v} is not a number", c.Param("id")));
		}

		_, err = database.DeleteProduct(id);
		if err != nil {
			return c.String(500, err.Error());
		}
		return c.String(200, "ok")
	});

	// todo convert into QueryParam and require auth
	e.POST("/api/v1/product/add", func (c echo.Context) error {
		payload := ProductBuffer {};
		err :=  (&echo.DefaultBinder{}).BindBody(c, &payload);

		if err != nil{
			return c.String(400, "Malfromated Product");
		}

		err = database.AddProduct(payload);
		if err != nil{
			return c.String(400, "failed to add product");
		}

		return c.JSON(200, payload);
	});

	e.POST("/api/v1/ticket/create_empty", func(c echo.Context) error {
		return c.String(404, "unimplemented!");
	})
	e.POST("/api/v1/ticket/create_width_sales", func(c echo.Context) error { return c.String(404, "unimplemented!"); })
	e.POST("/api/v1/ticket/:id/add_sale/:id", func(c echo.Context) error { return c.String(404, "unimplemented!"); })

	e.POST("/api/v1/sale/add/:name/:price/:count", func(c echo.Context) error { return c.String(404, "unimplemented!"); })

	// auth
	e.POST("/api/auth/v1/signin", func(c echo.Context) error {
		user := UserInfo{};
		err := (&echo.DefaultBinder{}).BindBody(c, &user);
		if err != nil{
			return c.String(500, "malfromed");
		}
	
		_, err = database.SignIn(user.Name, user.Password);
		if err != nil {
			return c.String(500, "could not signin");
		}
		return c.String(200, "signin succesful");
	});

	e.GET("/api/auth/v1/login", func(c echo.Context) error {
		user := UserInfo{};
		err := (&echo.DefaultBinder{}).BindBody(c, &user);
		if err != nil{
			return c.String(500, "malfromed");
		}
	
		_, err = database.LogIn(user.Name, user.Password);
		if err != nil{
			return c.String(400, "failed to login");
		}
		claim := &Claims{
			user.Name,
			jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		}}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.String(400, "failed to login");
		}
		return c.JSON(200, echo.Map{"token": t});
	});
	
}
