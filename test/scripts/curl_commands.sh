# get products
# curl \
#     --request GET \
#     "http://localhost:8080/api/v1/products"

# also get products
#
# curl --request GET "http://localhost:8080/api/v1/products" | jq

# get products by store
# curl \
#     --request GET \
#     "http://localhost:8080/api/v1/products?store=bcc%20tech"

# // add product
# curl \
#     --header "Content-Type: application/json" \
#     --request POST \
#     --data-raw '{"name":"ütü","price": 2000.0,"discount":10.0,"store":"bcc tech"}' \
#     "http://localhost:8080/api/v1/products"
#
# // update price
# curl \
#     --header "Content-Type: application/json" \
#     --request PUT \
#     "http://localhost:8080/api/v1/products/1?newPrice=4000.0" -v
#
#
# // delete product
# curl \
#     --header "Content-Type: application/json" \
#     --request DELETE \
#     "http://localhost:8080/api/v1/products/1" -v
