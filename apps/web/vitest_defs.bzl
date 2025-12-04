"""Helpers for running Vitest via Bazel."""

load("@npm//apps/web:vitest/package_json.bzl", vitest_bin = "bin")

def vitest_tests(name, config, deps, args = None, **kwargs):
    """Runs Vitest once (no watch) inside Bazel.

    Args:
        name: Name of the Bazel test target.
        config: Label of the `vitest.config` file.
        deps: Additional runtime dependencies, typically source filegroups.
        args: Extra CLI args (defaults to ["run"]).
    """

    user_env = kwargs.pop("env", {})
    merged_env = {
        "VITE_CJS_NODE_API": "1",
        "NODE_PATH": "$${RUNFILES_DIR}/_main/apps/web/node_modules",
    }
    merged_env.update(user_env)

    vitest_bin.vitest_test(
        name = name,
        args = args or ["run"],
        chdir = Label(config).package,
        data = [config] + deps + [
            ":node_modules",
            ":node_modules/jsdom",
            ":node_modules/vitest",
        ],
        env = merged_env,
        **kwargs
    )
