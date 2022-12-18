import math
import warnings
from collections import deque
from typing import List, Tuple
from helpers import *
from pptree import (  # INVESTIGATE: Consider using this library (in addition or instead) since it has custom typing: https://github.com/liwt31/print_tree
    Node,
    print_tree,
)
from simulator_types import RainTreeAnalytics, RainTreeQueueElement, RainTreeConfigs

warnings.filterwarnings("ignore", category=DeprecationWarning)

# A single RainTree propagation step
def propagate(
    elem: RainTreeQueueElement,
    analytics: RainTreeAnalytics,
    configs: RainTreeConfigs,
    queue: deque[RainTreeQueueElement],
) -> None:
    addr, sender, node, addr_book, depth = elem

    # Return if the addr book is empty
    if len(addr_book) == 0:
        return

    # Not a demote - real message over a network
    if addr != sender:
        analytics.msgs_rec_map[addr] += 1

    # If the theoretical depth was reached and no nodes are missing, return
    if len(analytics.nodes_missing) == 0:
        analytics.depth_reached_map[depth] += 1
        if depth >= configs.max_theoretical_depth:
            return

    # A network message was sent
    analytics.nodes_missing.discard(addr)
    analytics.nodes_reached.add(addr)

    # Configure who the current node should send messages to
    n = len(addr_book)
    i = addr_book.index(addr)
    t1 = (i + int(n * configs.t1_per)) % n
    t2 = (i + int(n * configs.t2_per)) % n
    s = (i + int(n * configs.shrinkage_per)) % n

    t1_addr = addr_book[t1]
    t2_addr = addr_book[t2]

    if t1_addr == t2_addr:
        t2_addr = None
    if t1_addr == addr:
        t1_addr = None

    def send(t: int, t_addr: str) -> None:
        analytics.msgs_sent += 1
        t_s = (t + int(n * configs.shrinkage_per)) % n
        t_book_s = shrink_list(addr_book.copy(), t, t_s)
        queue.append(
            (
                RainTreeQueueElement(
                    t_addr,
                    addr,
                    Node(t_addr, node),
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
        analytics.msgs_sent_map[addr] += 1
        print(f"Msg: {format_send_message(addr_book, i, t)}")

    # Send a message to the first target
    if t1_addr is not None:
        send(t1, t1_addr)

    # Send a message to the second target
    if t2_addr is not None:
        send(t2, t2_addr)

    # Demote - not incrementing `msg_send_counter` since it's not a send
    addr_book_s = shrink_list(addr_book, i, s)
    if len(addr_book_s) > 1:
        queue.append(
            (
                RainTreeQueueElement(
                    addr,
                    addr,
                    Node(addr, node),
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
                orig_addr,
                orig_addr,
                root_node,
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
    analytics: RainTreeAnalytics,
    raintreeConfigs: RainTreeConfigs,
) -> None:
    print("\n###################\n")
    print_tree(root_node, horizontal=False)
    print("\n###################\n")
    print(f"Coefficients used: t1: {raintreeConfigs.t1_per:.3f}, t2: {raintreeConfigs.t2_per:.3f}, shrinkage: {raintreeConfigs.shrinkage_per:.3f}")
    print(f"Num messages sent: {analytics.msgs_sent}")
    print(f"Num nodes reached: {len(analytics.nodes_reached)}/ {raintreeConfigs.num_nodes}")
    print(
        f"Messages received: {dict(dict(sorted(analytics.msgs_rec_map.items(), key=lambda item: -item[1])))}"
    )
    print(
        f"Messages sent: {dict(dict(sorted(analytics.msgs_sent_map.items(), key=lambda item: -item[1])))}"
    )
