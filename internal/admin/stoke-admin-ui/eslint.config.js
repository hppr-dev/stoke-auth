// eslint.config.mjs
import pluginVue from 'eslint-plugin-vue'
import {
  defineConfigWithVueTs,
  vueTsConfigs,
  configureVueProject,
} from '@vue/eslint-config-typescript'

// Optional: configure the Vue project to adjust the strictness of the rulesets or speed up linting.
configureVueProject({
  // Whether to parse TypeScript syntax in Vue templates.
  // Defaults to `true`.
  // Setting it to `false` could improve performance.
  // But TypeScript syntax in Vue templates will then lead to syntax errors.
  // Also, type-aware rules won't be applied to expressions in templates in that case.
  tsSyntaxInTemplates: true,

  // Whether to override some `no-unsafe-*` rules to avoid false positives on Vue component operations.
  // Defaults to `true`.
  // Usually you should keep this enabled,
  // but if you're using a metaframework or in a TSX-only project
  // where you're certain you won't operate on `.vue` components in a way that violates the rules,
  // and you want the strictest rules (e.g. when extending from `strictTypeChecked`),
  // you can set this to `false` to ensure the strictest rules apply to all files.
  allowComponentTypeUnsafety: true,

  // The root directory to resolve the `.vue` files.
  // Defaults to `process.cwd()`.
  // More info: <https://github.com/vuejs/eslint-plugin-vue/issues/1910#issuecomment-1819993961>
  // You may need to set this to the root directory of your project if you have a monorepo.
  // This is useful when you allow any other languages than `ts` in `.vue` files.
  // Our config helper would resolve and parse all the `.vue` files under `rootDir`,
  // and only apply the loosened rules to the files that do need them.
  rootDir: import.meta.dirname,
})

export default defineConfigWithVueTs(
  pluginVue.configs["flat/essential"],

  // We STRONGLY RECOMMEND you to start with `recommended` or `recommendedTypeChecked`.
  // But if you are determined to configure all rules by yourself,
  // you can start with `base`, and then turn on/off the rules you need.
  vueTsConfigs.base,
)
