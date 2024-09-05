# generator.py
"""
Script for creating a docker compose file with N amount of clients.
"""
import argparse


def generate_docker_compose_file(output_filename: str, num_clients: int) -> None:
    """
    Generates a docker compose file with N amount of clients.
    """
    with open("server-service.yaml", "r", encoding="utf-8") as server_file:
        server_string = server_file\
            .read()\
            .replace("TOTAL_AGENCIES=", f"TOTAL_AGENCIES={num_clients}")
    with open("network-tool.yaml", "r", encoding="utf-8") as network_file:
        network_string = network_file.read()
    with open(output_filename, "w", encoding="utf-8") as output:
        output.write(server_string)
        for i in range(num_clients):
            client_string = f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    entrypoint: /client
    environment:
        - CLI_ID={i}
        - FILE=/agency-{i}.csv
    networks:
        - testing_net
    depends_on:
        - server
    volumes:
      - ./client:/config
"""
            output.write(client_string)
        output.write(network_string)


# --------------------------------------------------------------------------------------------

parser = argparse.ArgumentParser(
    description="Generate a docker compose file with N amount of clients."
)
parser.add_argument("filename_output", type=str, help="Name of the output file.")
parser.add_argument("num_clients", type=int, help="Number of clients to generate.")
args = parser.parse_args()
print("\033[92mGenerating docker compose file...\033[0m")
try:
    generate_docker_compose_file(args.filename_output, args.num_clients)
except Exception as e:
    print("\033[91mAn error occurred while generating the docker compose file.\033[0m")
    raise e

with open(args.filename_output, "r", encoding="utf-8") as output_file:
    print(output_file.read())
print("\033[92mDocker compose file generated.\033[0m")
print(f"\033[92mFile: {args.filename_output}\033[0m")
print(f"\033[92mNumber of clients: {args.num_clients}\033[0m")
