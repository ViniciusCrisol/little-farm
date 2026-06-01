# Create Farm

curl -X POST http://localhost:8080/api/farms \
 -H "Content-Type: application/json" \
 -H "Accept-Encoding: gzip" \
 -d '{"map_width": 10, "map_height": 10}' | gunzip

# Find Farm

curl http://localhost:8080/api/farms/:farm_id \
 -H "Accept-Encoding: gzip" | gunzip

# Advance One Second

curl -X POST http://localhost:8080/api/farms/:farm_id/advance-one-second \
 -H "Accept-Encoding: gzip" | gunzip

# Plant Corn

curl -X POST http://localhost:8080/api/farms/:farm_id/dirt/:dirt_id/plant-corn \
 -H "Accept-Encoding: gzip" | gunzip

# Plant Wheat

curl -X POST http://localhost:8080/api/farms/:farm_id/dirt/:dirt_id/plant-wheat \
 -H "Accept-Encoding: gzip" | gunzip

# Harvest Corn

curl -X POST http://localhost:8080/api/farms/:farm_id/corn/:corn_id/harvest \
 -H "Accept-Encoding: gzip" | gunzip

# Harvest Grass

curl -X POST http://localhost:8080/api/farms/:farm_id/grass/:grass_id/harvest \
 -H "Accept-Encoding: gzip" | gunzip
