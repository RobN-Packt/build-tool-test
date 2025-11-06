const { execSync } = require("node:child_process");
const { join } = require("node:path");

module.exports = async function runExecutor(options, context) {
  const projectNode = context.projectGraph?.nodes?.[context.projectName];
  if (!projectNode) {
    throw new Error(`Nx project ${context.projectName} not found`);
  }
  const projectRoot = projectNode.data.root;
  const cwd = options.cwd ? join(projectRoot, options.cwd) : projectRoot;

  try {
    execSync(options.command, {
      cwd,
      stdio: "inherit",
      env: process.env,
    });
    return { success: true };
  } catch (error) {
    console.error("Go executor failed", error);
    return { success: false };
  }
};
