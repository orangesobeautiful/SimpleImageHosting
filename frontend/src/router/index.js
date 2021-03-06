import Vue from "vue";
import VueRouter from "vue-router";
import routes from "./routes";

// 處理 NavigationDuplicated 錯誤
// 在 router.push 重複路徑時，選擇 reload (重新整理)
const originalPush = VueRouter.prototype.push;
VueRouter.prototype.push = function push(location, onComplete, onAbort) {
  if (onComplete || onAbort)
    return originalPush.call(this, location, onComplete, onAbort);

  return originalPush.call(this, location).catch(failure => {
    if (
      VueRouter.isNavigationFailure(
        failure,
        VueRouter.NavigationFailureType.duplicated
      )
    ) {
      this.go(0);
    } else {
      throw failure;
    }
  });
};

Vue.use(VueRouter);

/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation;
 *
 * The function below can be async too; either use
 * async/await or return a Promise which resolves
 * with the Router instance.
 */

export default function(/* { store, ssrContext } */) {
  const Router = new VueRouter({
    scrollBehavior: () => ({ x: 0, y: 0 }),
    routes,

    // Leave these as they are and change in quasar.conf.js instead!
    // quasar.conf.js -> build -> vueRouterMode
    // quasar.conf.js -> build -> publicPath
    mode: process.env.VUE_ROUTER_MODE,
    base: process.env.VUE_ROUTER_BASE
  });

  return Router;
}
