package models

const JsonRest = `
		{
				"id": "ilf6q6610e4ellf",
				"name": "language_value",
				"type": "base",
				"system": false,
				"schema": [
						{
								"id": "w4gs2gol",
								"name": "language_key",
								"type": "relation",
								"system": false,
								"required": true,
								"unique": false,
								"options": {
										"collectionId": "5rxmpfs22bhonfh",
										"cascadeDelete": false,
										"minSelect": null,
										"maxSelect": 1,
										"displayFields": [
												"name"
										]
								}
						},
						{
								"id": "t4lju90v",
								"name": "iso",
								"type": "select",
								"system": false,
								"required": true,
								"unique": false,
								"options": {
										"maxSelect": 1,
										"values": [
												"de",
												"en"
										]
								}
						},
						{
								"id": "4bwzjixa",
								"name": "value",
								"type": "text",
								"system": false,
								"required": true,
								"unique": false,
								"options": {
										"min": null,
										"max": null,
										"pattern": ""
								}
						}
				],
				"listRule": "",
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
		},
		{
			"id": "yf7bx6hp7amhg5h",
			"created": "2023-03-13 08:43:14.709Z",
			"updated": "2023-03-13 08:43:14.709Z",
			"name": "page",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "rn2soypl",
					"name": "path",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "qiapmcpu",
					"name": "link",
					"type": "select",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"values": [
							"navbar",
							"burger",
							"footer",
							"hidden",
							"deactivated"
						]
					}
				},
				{
					"system": false,
					"id": "mvbtvhab",
					"name": "iso",
					"type": "select",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"values": [
							"de",
							"en"
						]
					}
				},
				{
					"system": false,
					"id": "ctrtizzo",
					"name": "title",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "npptexqi",
					"name": "content",
					"type": "editor",
					"required": true,
					"unique": false,
					"options": {}
				}
			],
			"listRule": "",
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}
`
