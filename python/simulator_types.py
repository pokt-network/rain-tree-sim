import math
from collections import defaultdict
from dataclasses import dataclass
from typing import List
from pptree import Node

# An object containing the configurations of how RainTree should propagate messages
@dataclass
class RainTreeConfigs:
    num_nodes: int
    t1_per: float # % of location in addr_book for the node's 1st target (whom to send a message to)
    t2_per: float # % of location in addr_book for the node's 1st target (whom to send a message to)
    shrinkage_per: float  # addr book shrinkage coefficient (i.e. % of addr book that should be kept for the next level)
    max_theoretical_depth: int
    addr_book: List[str]

    def __init__(
        self,
        num_nodes: int,
        t1_per: float,
        t2_per: float,
        shrinkage_per: float
    ):
        self.num_nodes = num_nodes
        self.t1_per = t1_per
        self.t2_per = t2_per
        self.shrinkage_per = shrinkage_per
        self.max_theoretical_depth = RainTreeConfigs._get_max_theoretical_depth(num_nodes)
        self.addr_book = RainTreeConfigs._generate_addr_book(num_nodes)


    def _get_max_theoretical_depth(num_nodes: int) -> int:
        return math.log(
            num_nodes, 3 # 3 comes from the fact that we use a ternary tree
        )

    def _generate_addr_book(num_nodes: int) -> List[str]:
        return sorted(
            [f"node_{i+1}" for i in range(num_nodes)], key=lambda x: int(x.split("_")[1])
        )

# An object containing a list of global counters intended to be propagated
# during RainTree broadcast to collect analytics on its performance.
@dataclass
class RainTreeAnalytics:
    msgs_sent: int # Total num of messages sent by RainTree propagating
    nodes_reached: set[str] # Nodes reached by current RainTree propagating
    nodes_missing: set[str] # Nodes not yet reached by current RainTree propagating
    msgs_rec_map: defaultdict[str, int] # Num messages received by per node by addr
    msgs_sent_map: defaultdict[str, int] # Num messages sent by node per addr
    depth_reached_map: defaultdict[str, int] # Max depth reached by node addr

    def __init__(
        self,
        nodes: List[str],
    ):
        self.msgs_sent = 0
        self.nodes_reached = set()
        self.nodes_missing = set(nodes)
        self.msgs_rec_map = defaultdict(int)
        self.msgs_sent_map = defaultdict(int)
        self.depth_reached_map = defaultdict(int)


# The python simulator of RainTree uses BFS (Breadth First Search) implemented via a FIFO queue
# where each element in the queue is an instance of this class
@dataclass
class RainTreeQueueElement:
    addr: str # The addr of the node handling the current message/propagation
    sender: str # sender addr (who sent the message to addr)
    node: Node # current node correspond to addr
    addr_book: List[str] # The addr book of the current node (i.e. list of sorted addresses)
    depth: int # the current depth in the propagations

    def __iter__(self):
        return iter(
            (
                self.addr,
                self.sender,
                self.node,
                self.addr_book,
                self.depth,
            )
        )
