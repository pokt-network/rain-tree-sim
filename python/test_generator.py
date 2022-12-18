import sys

import stringcase
from num2words import num2words
from pptree import Node, print_tree
from simulator import RainTreeAnalytics
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
    raintreeConfigs: RainTreeConfigs,
    analytics: RainTreeAnalytics,
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
