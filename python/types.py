from collections import defaultdict, deque
from dataclasses import dataclass


@dataclass
class Counters:
    msgs_sent: int  # Total num of messages sent by RainTree propagating
    nodes_reached: set[str]  # Nodes reached by current RainTree propagating
    nodes_missing: set[str]  # Nodes not yet reached by current RainTree propagating
    msgs_rec_map: defaultdict[str, int]  # Num messages received by node addr
    msgs_sent_map: defaultdict[str, int]  # Num messages sent by node addr
    depth_reached_map: defaultdict[str, int]  # Max depth reached by node addr
    # Theoretical max depth, used to end propagating early
    max_theoretical_depth: int

    def __init__(self, nodes: List[str], max_allowed_depth: int):
        self.msgs_sent = 0
        self.nodes_reached = set()
        self.nodes_missing = set(nodes)
        self.msgs_rec_map = defaultdict(int)
        self.msgs_sent_map = defaultdict(int)
        self.depth_reached_map = defaultdict(int)
        self.max_theoretical_depth = max_allowed_depth


@dataclass
class PropagationQueueElement:
    addr: str  # The addr of the node propagating the message
    addr_book: List[str]  # the addr's current addr book
    depth: int  # the current depth in the propagations
    t1: float  # % of location in addr book for target 1
    t2: float  # % of location in addr book for target 1
    shrinkage: float  # addr book shrinkage coefficient
    node: Node  # current node
    sender: str  # sender addr (who sent the message to addr)

    def __iter__(self):
        return iter(
            (
                self.addr,
                self.addr_book,
                self.depth,
                self.t1,
                self.t2,
                self.shrinkage,
                self.node,
                self.sender,
            )
        )
