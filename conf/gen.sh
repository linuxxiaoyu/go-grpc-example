openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650

openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 7200 -key ca.key -out ca.pem

openssl req -new -key server.key  -out server.csr
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server2.pem

openssl ecparam -genkey -name secp384r1 -out client.key
openssl req -new -key client.key -out client.csr