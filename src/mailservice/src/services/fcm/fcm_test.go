package fcm

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var KEY FB_PRIVATE_KEY = `-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDBo/jKSRlvBS5R\nh/r9CRbOXI+fp0XvEy3gGO2mnCX/NYkLjv3k45x6xvGHLM3UNvt18HU33PishYg/\n7Er86WXBHl2EapkPBL6kieFoEQkU4LAf3PtPtrriQOjvn1s7HIfyCig7/xC8Dwfl\np/0CNhZ6RDkNAy74wVzGkKZ0dajnksz1uFEcunQsUZmj80OOy3N4dxk4k1oMGcPO\n3l+voSrEVVXk5jmR3ciS5pdrrn1mZQNHTp98W01Cng1v+rF81a9QTwkp1LsFqBKD\nFVSLi7fRUiob2oAEKOh8X8ytqhJPDVa5yEfOdw4vmp6Msz6FCUzWeOLBt3VyfyHu\nG/KR5GJRAgMBAAECggEAGjqIn3XJUSVlgbumfpG1mhwlhB2XNmvlod4eipvJ9cid\nmIg00cUW0/aQjpu+AYm1A+OfLQLsWAn6S5ZJDfrbQo5HYoFB3CvrWsQmWP89uKs6\nkAZRsBlzNORP6O0v4VDbBSjlDENfU+nBSxU3Cw6iess04xNUUHN4ipjbQxkQ2NTo\nIw1k+7XveDogZegFwusGyV8CrtrjVNHtQXkkMo0Ub5j7e3DylokzVeJg5zZfgI24\nZqk0f5BPahr2X/fZ7fIJqCHZjHOsnxBNoo1Kz2aoztipc6z1TSrW8qAO19J5RfrI\nACndFvVnw59YyGHGChv9tJ3ziNvHxZtasE4i+sYUlQKBgQDubeTPOPtxkB/EsEbK\nQ3xR+kdwLyPtglpWJcAH/RyUC6gHjodJj6yYx43MaREsulwTVPy5yT4eZyr44euP\nILK7yPJ75CECsyZ09Zh/qB4xn5PxphXRPCpapXqboKr9IuvBJzRtUxozYszLjSgD\nUpkT7okkRIyxf4Nb5LutCJKdZwKBgQDP6Ry6PLcpMQNpnR32khXHiLBzqRkTb3CG\nGqxhzLlkfnCZOAyS3/eDZRfw8dnxKSi3/R3s8d9LzJY4fWJHXZu1uQko7yarHy0D\nAs4Wl2y177fPTgu2ej3Sk8JIjsLHS5BLzhqx/2Z3Rv1la2JTdCG6ohlwt3+XH3Lt\nEKX+vYP3hwKBgACKxWtnMMMoVbonwHFzR9QT4pexs741fqkVeuNJwwffIumpfEtB\nhV3vjjX5wy072zu8BLsTZw3ApEtekB+KLn3YzhxT/3M3Hw5DBK69nhv0xexVuVT5\ncwsztxyld94Nd0XAJhFdkACv59FKp92iEXEHKM6pTTyWEqFh2r9g9pxfAoGAZxGQ\nPT8mKdRzdaL/HKI1C9LWbrAQj1L6fHCyrlUYPxpzZXGkwhcnk8rFAJxUx7n4xqVD\ndZg+c0w72EtIMkrUi1Tslo9gIwr0fH6ifg6ZROROwgVVxyN4jHDVqrSjGLt8EChf\nkYgkWtMlgWanuuliYyxC4l8FcHyVs7JCKDP5PPcCgYAUJNWn6n6A/AJ+LnSLcV3N\nQbdgXucJdG6fLneQ2Dnan+r2fMIyunHN2Mi6wGphGNBAW3cnNHnH3XqWGIs1qCOw\nf70tuSxYjLkBYxwhgqaO8++DXI//HARiWb7g7S40Vpe7YSWAjzsRX36CtIiCL0Hu\nx5jrBN9J7FmMQzi7Iw5A1w==\n-----END PRIVATE KEY-----\n`
var (
	MY_ANDROID = "cy8bB1XgQl6rgTJTZEVOSg:APA91bFfTsuqjJfQjuJi4MtTkTIlhLu-SYsN8DaIN5vlxFOc4rXyWaL4QqWP4HDJE23UvH1GIuBb8dCez05y8usH4GZSWklDFptYDh9sLlIq2__I7wh7RMuiYO1sc3NwcezbV0mXTcJL"
	MY_VIRTUAL = "fc2M7AElRxiXH2Ultqs8Op:APA91bEnlnGenWTMf2Rk9y6rKj2Lz0f8vfPddwI3O5cvOgw4nGs1czs73i_7gEn72VleVQC1MUm5MG1sMxFVc-ABvGy1GnD34oDdvTLTxPEQuGW4uXmRAfqC-P8oyQNMyhlZNoNWbjMV"
)

func TestSendNotification(t *testing.T) {
	ctx := context.Background()
	client, err := NewService(logrus.New(), KEY)
	assert.Nil(t, err)
	id := uuid.New()
	err = client.SendTo(ctx, []string{MY_VIRTUAL}, "Title", "Anygonow Test message", map[string]string{
		"id":         id.String(),
		"seq":        fmt.Sprint(id.ClockSequence()),
		"type":       "reject-notification",
		"businessId": "006381d9-a8a2-479d-a507-84e6d3cf642c",
	})
	assert.Nil(t, err)

}
