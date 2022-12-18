from typing import Dict, List


# Returns a subset of the list with items from `i1` to `i2` in `l` wrapping around if needed
def shrink_list(l: List[str], i1: int, i2: int) -> List[str]:
    if i1 <= i2:
        return l[i1:i2]
    return l[i1:] + l[:i2]


# Formats a string showing who the sender and receivers were
def format_send_message(l: List[str], self_index: int, target: int) -> str:
    s = "[ "
    for idx, n in enumerate(l):
        if idx == self_index:
            s += f"({n}), "
        elif idx == target:
            s += f"**{n}**, "
        else:
            s += f"{n}, "
    return f"{s[:-2]} ]"


# Sum the values from two dictionaries where the keys overlap
def agg_dicts(d1: Dict[str, int], d2: Dict[str, int]) -> Dict[str, int]:
    return {k: d1.get(k, 0) + d2.get(k, 0) for k in set(d1) | set(d2)}
