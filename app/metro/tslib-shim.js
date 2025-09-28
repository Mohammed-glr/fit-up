// Resolve the actual CommonJS build of tslib once and re-export it with an explicit default
// to satisfy Metro's interop when packages import from "tslib/modules/index.js".
const tslib = require('../node_modules/tslib/tslib.js');

module.exports = {
  __esModule: true,
  default: tslib,
  ...tslib,
};
