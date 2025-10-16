module.exports = {
  extends: ["../../config/eslint-base.cjs"],
  parserOptions: {
    project: ["./tsconfig.json", "./tsconfig.eslint.json"]
  }
};
