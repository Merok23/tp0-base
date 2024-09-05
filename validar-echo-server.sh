docker network create tp0_net
docker network connect tp0_net server
input="hola"
response=$(echo "$input" | docker run --network tp0_net -i client:latest sh -c 'nc server 12345')
echo "Response: $response"
if [ "$input" = "$response" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi
docker network disconnect tp0_net server
docker network rm tp0_net