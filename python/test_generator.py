import argparse
import sys

import stringcase
from num2words import num2words
from pptree import Node, print_tree
from simulator import RainTreeAnalytics, display_simulation_results, simulate
from simulator_types import RainTreeConfigs

TEST_FORMAT = """
func TestRainTreeComplete{0}Nodes(t *testing.T) {{
    originatorNode := validatorId({1})
    var expectedCalls = TestNetworkSimulationConfig{{
{2}
    }}
	testRainTreeCalls(t, originatorNode, expectedCalls)
}}
"""


def prepare_test(
    root_node: Node,
    analytics: RainTreeAnalytics,
    raintreeConfigs: RainTreeConfigs,
    filename: str,
) -> None:
    test_generator = {}
    test = ""

    for k, _ in analytics.msgs_rec_map.items():
        test_generator[k] = (
            analytics.msgs_rec_map[k],
            analytics.msgs_sent_map[k],
        )

    originator_i = -1
    for i in range(raintreeConfigs.num_nodes):
        k = f"val_{i+1}"
        if k == root_node.name:
            originator_i = i + 1
            test += f"        originatorNode:"
        else:
            test += f"        validatorId({i+1}):"
        read = analytics.msgs_rec_map[k]
        write = analytics.msgs_sent_map[k]
        test += f"  {{{read}, {write}}},\n"

    num_nodes_words = stringcase.camelcase(
        num2words(raintreeConfigs.num_nodes).replace("-", "_")
    ).capitalize()
    go_test = TEST_FORMAT.format(num_nodes_words, originator_i, test)

    with open(filename, "w") as sys.stdout:
        print_tree(root_node, horizontal=False)
        print(go_test)


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
    display_simulation_results(root_node, analytics, raintreeConfigs)

    # Prepare Test
    prepare_test(root_node, analytics, raintreeConfigs, args.output_file)


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
