const routes = [
  {
    path: "/",
    component: () => import("layouts/MainLayout.vue"),
    children: [
      { path: "", component: () => import("pages/Index.vue") },
      { path: "upload", component: () => import("pages/Upload.vue") },
      {
        path: "image/:imgHashID",
        component: () => import("pages/Image.vue")
      },
      {
        path: "user/:id/",
        component: () => import("pages/User.vue"),
        children: [
          {
            path: "images",
            component: () => import("pages/User/Images.vue")
          },
          {
            path: "about",
            component: () => import("pages/User/About.vue")
          }
        ]
      }
    ]
  },
  {
    path: "/dashboard",
    component: () => import("layouts/MainLayout.vue"),
    children: [
      {
        path: "",
        component: () => import("pages/dashboard/dashboard.vue")
      },
      {
        path: "settings",
        component: () => import("pages/dashboard/settings.vue")
      }
    ]
  },
  {
    path: "/account-activate/:token",
    component: () => import("layouts/EmptyLayout.vue"),
    children: [
      {
        path: "",
        component: () => import("pages/AccountActivate.vue")
      }
    ]
  },
  {
    path: "/register",
    component: () => import("layouts/EmptyLayout.vue"),
    children: [
      {
        path: "",
        component: () => import("pages/Register.vue")
      }
    ]
  },
  {
    path: "/signin",
    component: () => import("layouts/EmptyLayout.vue"),
    children: [
      {
        path: "",
        component: () => import("pages/Signin.vue")
      }
    ]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: "*",
    component: () => import("pages/Error404.vue")
  }
];

export default routes;
