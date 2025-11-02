module.exports = {
  extends: ["../../config/eslint-base.cjs"],
  parserOptions: {
    project: ["./tsconfig.eslint.json"],
    tsconfigRootDir: __dirname
  }
};
