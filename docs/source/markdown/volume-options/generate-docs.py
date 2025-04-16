import json
import sys


def alternative_names(parsed_json, opt):
    alternative_names = parsed_json["options"][opt]["alternative_name"]
    if len(alternative_names) == 0:
        return ""
    else:
        return "(" + ",".join(alternative_names) + ")"


def argcount_compatibility(parsed_json, opt):
    str = "| " + opt + " | "
    for i in ["argcount_0", "argcount_1"]:
        if parsed_json["options"][opt][i]:
            str += "x"
        str += " | "
    return str


def option_description(parsed_json, opt):
    supported_markdown_formatted = [
        "_" + item + "_" for item in sorted(parsed_json["options"][opt]["supported"])
    ]
    str = (
        "#### "
        + opt
        + alternative_names(parsed_json, opt)
        + "\n\n"
        + parsed_json["options"][opt]["docs"]
        + "\n"
        + "Available for types: "
        + ", ".join(supported_markdown_formatted)
    )
    return str


def mounttype_description(parsed_json, mounttype):
    supported_options = []
    for opt in sorted(parsed_json["options"]):
        if mounttype in parsed_json["options"][opt]["supported"]:
            supported_options.append("_" + opt + "_")
    str = (
        "#### "
        + mounttype
        + "\n\n"
        + parsed_json["types"][mounttype]["description"]
        + "\n"
        + "Available options: "
        + ", ".join(supported_options)
    )
    return str


def supported_table(parsed_json):
    volumetypes = sorted(parsed_json["types"].keys())
    header_text = "| |" + " | ".join(volumetypes) + " |"
    header_line = "| --- |" + " | ".join(["---"] * len(volumetypes)) + " |"
    rows = [
        "| "
        + opt
        + alternative_names(parsed_json, opt)
        + " | "
        + " | ".join(
            "x" if volumetype in parsed_json["options"][opt]["supported"] else " "
            for volumetype in volumetypes
        )
        + " |"
        for opt in sorted(parsed_json["options"])
    ]
    return "\n".join([header_text, header_line] + rows)


def argcount_table(parsed_json):
    keys = sorted(parsed_json["types"].keys())
    header_text = "| | 0 arg allowed | 1 arg allowed |"
    header_line = "| - | - | - |"
    rows = [
        argcount_compatibility(parsed_json, opt)
        for opt in sorted(parsed_json["options"])
    ]
    return "\n".join([header_text, header_line] + rows)


def option_list(parsed_json):
    options = [
        option_description(parsed_json, opt)
        for opt in sorted(parsed_json["options"].keys())
    ]
    res = "\n\n### Options for `--mount`:\n\n" + "\n\n".join(options)
    return res


def type_list(parsed_json):
    mounttypes = sorted(parsed_json["types"].keys())
    descriptions = [
        mounttype_description(parsed_json, mounttype) for mounttype in mounttypes
    ]
    res = "\n\n### Types for `--mount`:\n\n" + "\n\n".join(descriptions)
    return res


def main() -> int:
    try:
        data = sys.stdin.read()
        parsed_json = json.loads(data)
        print(supported_table(parsed_json))
        print("\n")
        print(argcount_table(parsed_json))
        print(type_list(parsed_json))
        print(option_list(parsed_json))
        return 0

    except json.JSONDecodeError as e:
        print("error: invalid json. " + str(e), file=sys.stderr)
        return 1

    except json.UnicodeDecodeError as e:
        print("error: invalid unicode. " + str(e), file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
