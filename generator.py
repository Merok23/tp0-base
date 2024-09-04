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
        server_string = server_file.read()
    with open("network-tool.yaml", "r", encoding="utf-8") as network_file:
        network_string = network_file.read()
    print("\033[92mGenerating docker compose file...\033[0m")
    with open(output_filename, "w", encoding="utf-8") as output_file:
        output_file.write(server_string)
        for i in range(1, num_clients + 1):
            client_string = f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    entrypoint: /client
    environment:
        - CLI_LOG_LEVEL=DEBUG
        - CLI_ID={i}
    networks:
        - testing_net
    depends_on:
        - server
"""
            output_file.write(client_string)
        output_file.write(network_string)
    with open(output_filename, "r", encoding="utf-8") as output_file:
        print(output_file.read())
    print("\033[92mDocker compose file generated.\033[0m")
    print(f"\033[92mFile: {output_filename}\033[0m")
    print(f"\033[92mNumber of clients: {num_clients}\033[0m")


# --------------------------------------------------------------------------------------------

parser = argparse.ArgumentParser(
    description="Generate a docker compose file with N amount of clients."
)
parser.add_argument("filename_output", type=str, help="Name of the output file.")
parser.add_argument("num_clients", type=int, help="Number of clients to generate.")
args = parser.parse_args()
generate_docker_compose_file(args.filename_output, args.num_clients)
