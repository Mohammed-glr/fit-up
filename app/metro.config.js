const path = require('path');
const { getDefaultConfig } = require('expo/metro-config');
const { resolve } = require('metro-resolver');

const projectRoot = __dirname;
const config = getDefaultConfig(projectRoot);

const tslibShimPath = path.join(projectRoot, 'metro', 'tslib-shim.js');

const baseResolveRequest = config.resolver?.resolveRequest;

config.resolver = {
  ...config.resolver,
  resolveRequest(context, moduleName, platform) {
    if (
      moduleName === 'tslib' ||
      moduleName === 'tslib/modules/index.js' ||
      moduleName === 'tslib/modules/index'
    ) {
      return {
        type: 'sourceFile',
        filePath: tslibShimPath,
      };
    }

    if (typeof baseResolveRequest === 'function') {
      return baseResolveRequest(context, moduleName, platform);
    }

    return resolve(context, moduleName, platform);
  },
  extraNodeModules: {
    ...(config.resolver?.extraNodeModules ?? {}),
    tslib: tslibShimPath,
  },
};

module.exports = config;
