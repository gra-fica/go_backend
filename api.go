package main

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

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

	// auth apis
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

	e.POST("/api/v1/auth/v1/signin", func(c echo.Context) error {
		return nil;
	});
}
