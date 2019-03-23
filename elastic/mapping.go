package elastic

const ProducerMapping = `
{
	"settings":{
		"number_of_shards": 2,
		"number_of_replicas": 0
	},
	"mappings":{
		"producer":{
			"properties":{
				"account":{
					"type":"keyword"
				},
				"name":{
					"type":"text"
				},
				"category":{
					"type":"text",
					"store": true,
					"fielddata": false
				},
				"startdate":{
					"type":"date",
					"format": "yyyy-MM-dd"
				},
				"enddate":{
					"type":"date",
					"format": "yyyy-MM-dd"
				}, 
				"location":{
					"type":"geo_point"					 
				}, 
				"quantity":{
					"type":"float"					 
				},
				"icon":{
					"type":"text"
				},
				"client":{
					"type":"text"
				}
			}
		}
	}
}
`

const ConsumerMapping = `
{
	"settings":{
		"number_of_shards": 2,
		"number_of_replicas": 0
	},
	"mappings":{
		"consumer":{
			"properties":{
				"account":{
					"type":"keyword"
				},
				"name":{
					"type":"text"
				},
				"category":{
					"type":"text",
					"store": true,
					"fielddata": false
				},
				"startdate":{
					"type":"date",
					"format": "yyyy-MM-dd"
				},
				"enddate":{
					"type":"date",
					"format": "yyyy-MM-dd"
				}, 
				"location":{
					"type":"geo_point"					 
				}, 
				"quantity":{
					"type":"float"					 
				},
				"icon":{
					"type":"text"
				}  
			}
		}
	}
}
`
