import math
import warnings
from collections import deque
from typing import List, Tuple

from helpers import *
from pptree import (  # INVESTIGATE: Consider using this library (in addition or instead) since it has custom typing: https://github.com/liwt31/print_tree
    Node,
    print_tree,
)
from simulator_types import RainTreeAnalytics, RainTreeConfigs, RainTreeQueueElement

warnings.filterwarnings("ignore", category=DeprecationWarning)

# A single RainTree propagation step
def propagate(
    elem: RainTreeQueueElement,
    analytics: RainTreeAnalytics,
    configs: RainTreeConfigs,
    queue: deque[RainTreeQueueElement],
) -> None:
    node, sender, addr_book, depth = elem

    # Return if the addr book is empty
    if len(addr_book) == 0:
        return

    # Not a demote - real message over a network
    if node.name != sender:
        analytics.msgs_rec_map[node.name] += 1

    # A network message was received
    analytics.nodes_missing.discard(node.name)
    analytics.nodes_reached.add(node.name)

    # If the theoretical depth was reached and no nodes are missing, return
    if len(analytics.nodes_missing) == 0:
        analytics.depth_reached_map[depth] += 1
        if depth >= configs.max_theoretical_depth:
            return

    # Define function to perform a network send
    def send(t: int, t_addr: str) -> None:
        analytics.msgs_sent += 1
        t_s = (t + int(num_nodes * configs.shrinkage_per)) % num_nodes
        t_book_s = shrink_list(addr_book.copy(), t, t_s)
        queue.append(
            (
                RainTreeQueueElement(
                    Node(t_addr, node),
                    node.name,
                    t_book_s,
                    depth + 1,
                ),
                analytics,
                configs,
                queue,
            )
        )

        analytics.nodes_missing.discard(t_addr)
        analytics.nodes_reached.add(t_addr)
        analytics.msgs_sent_map[node.name] += 1
        print(f"Msg: {format_send_message(addr_book, node_idx, t)}")

    # Configure who the current node should send messages to
    num_nodes = len(
        addr_book
    )  # not using configs.num_nodes because we need the length of the current addressb ook
    node_idx = addr_book.index(node.name)
    t1_idx = (node_idx + int(num_nodes * configs.t1_per)) % num_nodes
    t2_idx = (node_idx + int(num_nodes * configs.t2_per)) % num_nodes

    t1_addr = addr_book[t1_idx]
    t2_addr = addr_book[t2_idx]

    if t1_addr == t2_addr:
        t2_addr = None
    if t1_addr == node.name:
        t1_addr = None

    # Send a message to the target 1
    if t1_addr is not None:
        send(t1_idx, t1_addr)

    # Send a message to the target 2
    if t2_addr is not None:
        send(t2_idx, t2_addr)

    # Demote - not incrementing `msg_send_counter` since it's not a send
    s = (
        node_idx + int(num_nodes * configs.shrinkage_per)
    ) % num_nodes  # shrunk_addr_book_end_index
    addr_book_s = shrink_list(addr_book, node_idx, s)
    if len(addr_book_s) > 1:
        queue.append(
            (
                RainTreeQueueElement(
                    Node(node.name, node),
                    node.name,
                    addr_book_s,
                    depth + 1,
                ),
                analytics,
                configs,
                queue,
            )
        )


# A single RainTree Simulation
def simulate(
    orig_addr: str,
    raintreeConfigs: RainTreeConfigs,
) -> Tuple[Node, RainTreeAnalytics]:
    # Configure Simulation
    queue = deque()
    analytics = RainTreeAnalytics(raintreeConfigs.addr_book)

    # Prepare Simulation
    root_node = Node(orig_addr)
    queue.append(
        (
            RainTreeQueueElement(
                root_node,
                orig_addr,
                raintreeConfigs.addr_book,
                0,
            ),
            analytics,
            raintreeConfigs,
            queue,
        )
    )

    # Run Simulation to completion
    while len(queue) > 0:
        propagate(*queue.popleft())

    return root_node, analytics


# TODO(olshansky): Make sure these are colored in the output in terminal for easier readability
def display_simulation_results(
    root_node: Node,
    raintreeConfigs: RainTreeConfigs,
    analytics: RainTreeAnalytics,
) -> None:
    print("\n###################\n")
    print_tree(root_node, horizontal=False)
    print("\n###################\n")
    print(
        f"Coefficients used: t1: {raintreeConfigs.t1_per:.3f}, t2: {raintreeConfigs.t2_per:.3f}, shrinkage: {raintreeConfigs.shrinkage_per:.3f}"
    )
    print(f"Num messages sent: {analytics.msgs_sent}")
    print(
        f"Num nodes reached: {len(analytics.nodes_reached)}/ {raintreeConfigs.num_nodes}"
    )
    print(
        f"Messages received: {dict(dict(sorted(analytics.msgs_rec_map.items(), key=lambda item: -item[1])))}"
    )
    print(
        f"Messages sent: {dict(dict(sorted(analytics.msgs_sent_map.items(), key=lambda item: -item[1])))}"
    )
