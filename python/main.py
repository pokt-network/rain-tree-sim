import argparse
from simulator import display_simulation_results, simulate
from simulator_types import RainTreeConfigs
from test_generator import prepare_test

def main(args):
    # Simulation Parameters
    raintreeConfigs = RainTreeConfigs(
        args.num_nodes,
        args.t1_per,
        args.t2_per,
        args.shrinkage_per,
    )
    orig_addr = raintreeConfigs.addr_book[0]

    # Run Simulation
    root_node, analytics = simulate(orig_addr, raintreeConfigs)

    # Print Results
    display_simulation_results(root_node, raintreeConfigs, analytics)

    # Prepare Test
    prepare_test(root_node, raintreeConfigs, analytics, args.output_file)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--num_nodes",
        dest="num_nodes",
        type=int,
        default=42,
        help="# of nodes used to simulated RainTree broadcast",
    )
    parser.add_argument(
        "--t1_per",
        dest="t1_per",
        type=float,
        default=1/3,
        help="% target of first message in the sorted address book",
    )
    parser.add_argument(
        "--t2_per",
        dest="t2_per",
        type=float,
        default=2/3,
        help="% target of first message in the sorted address book",
    )
    parser.add_argument(
        "--shrinkage_per",
        dest="shrinkage_per",
        type=float,
        default=2/3,
        help="% shrinkage of addr book with each decreased level",
    )
    parser.add_argument(
        "--output_file",
        dest="output_file",
        type=str,
        default="raintree_single_test.go",
        help="Output file where the generated Golang test should be written to",
    )
    args = parser.parse_args()
    main(args)
