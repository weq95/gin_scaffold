package public

import (
	"encoding/json"
	"fmt"
)

var tree = `[
    {
      "id": "1",
      "level": "1",
      "pid": "0",
      "name": "服装",
      "icon": "/static/uploadtest/categoryIcon/4a14b14397e117b551b7bb764925cf6c.jpg",
      "ctime": "0",
      "uptime": "1618042239",
      "sort": "0"
    },
    {
      "id": "2",
      "level": "2",
      "pid": "1",
      "name": "上衣",
      "icon": "/static/uploadtest/categoryIcon/be1f9cc0f850b2a3b514817e72d7db50.jpg",
      "ctime": "0",
      "uptime": "1613705080",
      "sort": "0"
    },
    {
      "id": "4",
      "level": "3",
      "pid": "2",
      "name": "毛衣",
      "icon": "/static/uploadtest/categoryIcon/7821109bba4cd7d017e027e1d1c9c2f4.jpg",
      "ctime": "0",
      "uptime": "1623808882",
      "sort": "1"
    },
    {
      "id": "5",
      "level": "3",
      "pid": "2",
      "name": "打底衫",
      "icon": "",
      "ctime": "0",
      "uptime": "0",
      "sort": "0"
    },
    {
      "id": "6",
      "level": "1",
      "pid": "0",
      "name": "美妆",
      "icon": "/static/uploadtest/categoryIcon/0322308aa64c45e1edbb5847e71456f0.jpg",
      "ctime": "0",
      "uptime": "1612495801",
      "sort": "0"
    },
    {
      "id": "7",
      "level": "2",
      "pid": "6",
      "name": "护肤品",
      "icon": "",
      "ctime": "0",
      "uptime": "1609833953",
      "sort": "0"
    },
    {
      "id": "8",
      "level": "2",
      "pid": "6",
      "name": "彩妆",
      "icon": "",
      "ctime": "0",
      "uptime": "0",
      "sort": "0"
    },
    {
      "id": "9",
      "level": "3",
      "pid": "8",
      "name": "香水",
      "icon": "",
      "ctime": "0",
      "uptime": "0",
      "sort": "0"
    },
    {
      "id": "10",
      "level": "3",
      "pid": "8",
      "name": "BB霜",
      "icon": "",
      "ctime": "0",
      "uptime": "0",
      "sort": "0"
    },
    {
      "id": "11",
      "level": "1",
      "pid": "0",
      "name": "母婴玩具",
      "icon": "/static/uploadtest/categoryIcon/711fef18b05500827fb6ae1dcfd0ab51.jpg",
      "ctime": "1608599997",
      "uptime": "1608599997",
      "sort": "0"
    },
    {
      "id": "12",
      "level": "2",
      "pid": "11",
      "name": "车床用品",
      "icon": "/static/uploadtest/categoryIcon/c6bbefe05affb085552de3235a80e4d4.jpg",
      "ctime": "1608600262",
      "uptime": "1608600262",
      "sort": "0"
    },
    {
      "id": "15",
      "level": "2",
      "pid": "11",
      "name": "2222",
      "icon": "/static/uploadtest/categoryIcon/5c8b92e15ac54da99c828299435583bd.jpg",
      "ctime": "1609834124",
      "uptime": "1613704164",
      "sort": "0"
    },
    {
      "id": "16",
      "level": "1",
      "pid": "0",
      "name": "食品",
      "icon": "",
      "ctime": "1610440317",
      "uptime": "1610440317",
      "sort": "0"
    },
    {
      "id": "17",
      "level": "1",
      "pid": "0",
      "name": "生鲜水果",
      "icon": "/static/uploadtest/categoryIcon/ea363cdc9a9199fc0ea78f913f7e8c4e.jpg",
      "ctime": "1610594042",
      "uptime": "1610594419",
      "sort": "0"
    },
    {
      "id": "18",
      "level": "2",
      "pid": "17",
      "name": "新鲜蔬菜",
      "icon": "/static/uploadtest/categoryIcon/e42048a4fe78e01babb5bc8c4ef43617.jpg",
      "ctime": "1610594339",
      "uptime": "1610594339",
      "sort": "0"
    },
    {
      "id": "19",
      "level": "3",
      "pid": "18",
      "name": "土豆",
      "icon": "/static/uploadtest/categoryIcon/4d15d1a62424e276483ad194d0ff34b6.jpg",
      "ctime": "1610594380",
      "uptime": "1610594380",
      "sort": "4"
    },
    {
      "id": "22",
      "level": "1",
      "pid": "0",
      "name": "家电",
      "icon": "",
      "ctime": "1610597627",
      "uptime": "1610597627",
      "sort": "0"
    },
    {
      "id": "23",
      "level": "2",
      "pid": "22",
      "name": "空调",
      "icon": "/static/uploadtest/categoryIcon/96a3fb644d7e187e010e91f4a72ca8e8.jpg",
      "ctime": "1610611159",
      "uptime": "1610611159",
      "sort": "0"
    },
    {
      "id": "24",
      "level": "3",
      "pid": "23",
      "name": "智能",
      "icon": "/static/uploadtest/categoryIcon/9a64e60ad02a1d83c3fb420c0080bb1c.jpg",
      "ctime": "1610615168",
      "uptime": "1610615168",
      "sort": "0"
    },
    {
      "id": "25",
      "level": "1",
      "pid": "0",
      "name": "家具建材",
      "icon": "",
      "ctime": "1610637268",
      "uptime": "1623757686",
      "sort": "0"
    },
    {
      "id": "26",
      "level": "2",
      "pid": "29",
      "name": "面膜",
      "icon": "/static/uploadtest/categoryIcon/881cc326596655f7c57269105cdf6f6b.jpg",
      "ctime": "1611732896",
      "uptime": "1611733419",
      "sort": "0"
    },
    {
      "id": "27",
      "level": "1",
      "pid": "0",
      "name": "快乐源泉",
      "icon": "/static/uploadtest/categoryIcon/17ad6671b4781167b5bb5c5fd8840358.jpg",
      "ctime": "1611920182",
      "uptime": "1611920323",
      "sort": "0"
    },
    {
      "id": "28",
      "level": "2",
      "pid": "27",
      "name": "1213",
      "icon": "/static/uploadtest/categoryIcon/236f0e3d7b33cc7b92c8f7a3dd0ea4fd.jpg",
      "ctime": "1611920260",
      "uptime": "1611920332",
      "sort": "0"
    },
    {
      "id": "29",
      "level": "1",
      "pid": "0",
      "name": "汽车",
      "icon": "/static/uploadtest/categoryIcon/52b4073a753a40a486bfcad4c76b14ef.jpg",
      "ctime": "1612409592",
      "uptime": "1612409592",
      "sort": "0"
    },
    {
      "id": "30",
      "level": "2",
      "pid": "29",
      "name": "卡车",
      "icon": "/static/uploadtest/categoryIcon/9e8fcf8296e22aff16c5c793c0a2a9d1.jpg",
      "ctime": "1612409653",
      "uptime": "1612409653",
      "sort": "0"
    },
    {
      "id": "32",
      "level": "3",
      "pid": "2",
      "name": "23(⊙o⊙)…2423",
      "icon": "/static/uploadtest/categoryIcon/5d5c127250f1ba95909e546c28b510ed.jpg",
      "ctime": "1612494838",
      "uptime": "1612494838",
      "sort": "0"
    },
    {
      "id": "34",
      "level": "2",
      "pid": "25",
      "name": "3213",
      "icon": "/static/uploadtest/categoryIcon/adf20b238fa8bbdaba9f9ac8aa9824db.jpg",
      "ctime": "1623757856",
      "uptime": "1623757856",
      "sort": "0"
    },
    {
      "id": "35",
      "level": "2",
      "pid": "8",
      "name": "3213213",
      "icon": "/static/uploadtest/categoryIcon/3b822a38167e9ba3ad32fcd4e02ccf73.jpg",
      "ctime": "1623757871",
      "uptime": "1623757871",
      "sort": "321"
    },
    {
      "id": "36",
      "level": "3",
      "pid": "7",
      "name": "3213",
      "icon": "/static/uploadtest/categoryIcon/afe14c9f4a29bd78d00d2a08696861f5.jpg",
      "ctime": "1623757889",
      "uptime": "1623757889",
      "sort": "0"
    },
    {
      "id": "37",
      "level": "1",
      "pid": "0",
      "name": "312·12·1312",
      "icon": "/static/uploadtest/categoryIcon/e1aee1e56524f6ebacca7b208a187bd0.jpg",
      "ctime": "1623757906",
      "uptime": "1623757972",
      "sort": "321"
    },
    {
      "id": "38",
      "level": "2",
      "pid": "6",
      "name": "咋形成",
      "icon": "/static/uploadtest/categoryIcon/3b4dfe33a6935952197aecef95565a0e.jpg",
      "ctime": "1623758019",
      "uptime": "1623758019",
      "sort": "0"
    },
    {
      "id": "39",
      "level": "3",
      "pid": "18",
      "name": "番茄",
      "icon": "/static/uploadtest/categoryIcon/4d15d1a62424e276483ad194d0ff34b6.jpg",
      "ctime": "1610594380",
      "uptime": "1610594380",
      "sort": "1"
    },
    {
      "id": "40",
      "level": "3",
      "pid": "18",
      "name": "玉米",
      "icon": "/static/uploadtest/categoryIcon/4d15d1a62424e276483ad194d0ff34b6.jpg",
      "ctime": "1610594380",
      "uptime": "1610594380",
      "sort": "3"
    },
    {
      "id": "41",
      "level": "3",
      "pid": "18",
      "name": "南瓜",
      "icon": "/static/uploadtest/categoryIcon/4d15d1a62424e276483ad194d0ff34b6.jpg",
      "ctime": "1610594380",
      "uptime": "1610594380",
      "sort": "5"
    },
    {
      "id": "42",
      "level": "3",
      "pid": "18",
      "name": "黄瓜",
      "icon": "/static/uploadtest/categoryIcon/4d15d1a62424e276483ad194d0ff34b6.jpg",
      "ctime": "1610594380",
      "uptime": "1610594380",
      "sort": "2"
    }
  ]`

type Category struct {
	Id       string `json:"id"`
	Level    string `json:"level"`
	Pid      string `json:"pid"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Ctime    string `json:"ctime"`
	UpTime   string `json:"up_time"`
	Sort     string `json:"sort"`
	Children Trees  `json:"children"`
}

type Trees []*Category

func (tree Trees) ToTree() Trees {
	mi := make(map[string]*Category)
	for _, category := range tree {
		mi[category.Id] = category
	}

	var list = make([]*Category, 0)
	for _, category := range tree {
		if category.Pid == "0" {
			list = append(list, category)
			continue
		}

		if pItem, ok := mi[category.Pid]; ok {
			if pItem.Children == nil {
				pItem.Children = Trees{category}
				continue
			}

			pItem.Children = append(pItem.Children, category)
		}
	}

	return list
}

func main() {

	var category Trees
	err := json.Unmarshal([]byte(tree), &category)

	if err != nil {
		fmt.Println("category err:", err.Error())
		return
	}

	children := make([]*Category, 0)
	tree := make(Trees, len(category))
	for i, c := range category {
		c.Children = children
		tree[i] = c
	}

	newTree := tree.ToTree()

	treeByte, err := json.Marshal(newTree)
	if err != nil {
		fmt.Println("output err:", err.Error())
		return
	}

	fmt.Println(string(treeByte))
}
