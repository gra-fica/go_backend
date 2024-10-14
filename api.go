package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/labstack/echo-jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"

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

	e.GET("/api/v1/product/search", func (c echo.Context) error {
		name := c.QueryParam("name");
		count, count_err := strconv.Atoi( c.QueryParam("total") );

		products, err := database.SearchProduct(name, &ListAllProducts{}, &TokenFuzzy{});
		if count_err == nil{
			products = products[0:count];
			fmt.Println("invalid count");
		}

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
		if err != nil {
            fmt.Printf("%v", err);
			return c.String(404, fmt.Sprintf("\"%s\" not found", name));
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

	e.GET("/api/v1/debug/user/list", func(c echo.Context) error {
		rows, err := database.Query("GET-USERS");
		user := User{}
		if err != nil{
			return c.String(500, "no users");
		}
		for {
			err := rows.Scan(&user.ID, &user.Name, &user.Role, &user.Password);
			fmt.Printf("%d %s %s %s\n", user.ID, user.Name, user.Role, user.Password);
			if rows.Next() || err != nil{
				break;
			}
		}
		return c.String(200, "");
	});

	// need auth apis
	// todo convert into QueryParam and require auth
	g := e.Group("/api/v1");
	g.Use(echojwt.JWT([]byte("secret")));
	g.DELETE("/product/delete/:id", func (c echo.Context) error {
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
	g.POST("/product/add", func (c echo.Context) error {
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

	g.POST("/ticket/create_empty", func(c echo.Context) error {
		return c.String(404, "unimplemented!");
	})
	g.POST("/ticket/create_width_sales", func(c echo.Context) error { return c.String(404, "unimplemented!"); })
	g.POST("/ticket/:id/add_sale/:id", func(c echo.Context) error { return c.String(404, "unimplemented!"); })

	g.POST("/sale/add/:name/:price/:count", func(c echo.Context) error { return c.String(404, "unimplemented!"); })
	
}
