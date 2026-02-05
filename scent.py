from sniffer.api import *
import os
import termstyle
from subprocess import call

pass_fg_color = termstyle.green
pass_bg_color = termstyle.bg_default
fail_fg_color = termstyle.red
fail_bg_color = termstyle.bg_default

watch_paths = ["."]

ignore_dirs = []


def find_go_test_dirs():
    test_dirs = set()
    for root, dirs, files in os.walk("."):
        dirs[:] = [d for d in dirs if not d.startswith(".") and d not in ignore_dirs]
        for f in files:
            if f.endswith("_test.go"):
                test_dirs.add(root)
                break
    return test_dirs


@file_validator
def go_files(filename):
    path_parts = filename.split(os.sep)
    if any(ignored in path_parts for ignored in ignore_dirs):
        return False
    return filename.endswith(".go") and not os.path.basename(filename).startswith(".")


@runnable
def run_go_tests(*args):
    test_dirs = find_go_test_dirs()
    if not test_dirs:
        print(fail_fg_color("No test directories found"))
        return True

    all_passed = True
    for test_dir in sorted(test_dirs):
        print(f"\n{pass_fg_color('=' * 50)}")
        print(pass_fg_color(f"Running tests in: {test_dir}"))
        print(pass_fg_color("=" * 50))
        result = call(["go", "test", "-v", "./..."], cwd=test_dir)
        if result != 0:
            all_passed = False
            print(fail_fg_color(f"FAILED: {test_dir}"))
        else:
            print(pass_fg_color(f"PASSED: {test_dir}"))

    if all_passed:
        print(f"\n{pass_fg_color('All tests passed!')}")
    else:
        print(f"\n{fail_fg_color('Some tests failed!')}")

    return all_passed
