<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="toolbar">
      <div>
        <q-toolbar class="text-black">
          <q-toolbar-title>
            <q-btn flat @click="push.homePage()"> Simple Image Hosting </q-btn>
          </q-toolbar-title>

          <div class="row" v-if="loginInfoChecked">
            <div v-if="userID < 0">
              <q-btn @click="push.registerPage()">Register</q-btn>
              <q-btn @click="push.signinPage()">Signin</q-btn>
            </div>
            <div v-if="userID >= 0" class="row items-center">
              <div class="row">{{ showName }}</div>
              <q-avatar class="q-ml-sm">
                <img :src="avatar" />
                <q-menu :style="{ backgroundColor: '#eee', color: 'blue' }">
                  <q-list style="min-width: 100px">
                    <q-item clickable @click="push.userImagesPage(userID)">
                      <q-item-section>我的圖片</q-item-section>
                    </q-item>
                    <q-separator />
                    <q-item clickable @click="push.uploadPage()">
                      <q-item-section>上傳</q-item-section>
                    </q-item>
                    <q-separator />
                    <q-item
                      v-if="grade == 1"
                      clickable
                      @click="push.dashboardPage()"
                    >
                      <q-item-section>控制台</q-item-section>
                    </q-item>
                    <q-separator />
                    <q-item clickable @click="logout">
                      <q-item-section>登出</q-item-section>
                    </q-item>
                  </q-list>
                </q-menu>
              </q-avatar>
            </div>
          </div>
        </q-toolbar>
      </div>
    </q-header>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<style lang="sass" scoped>
.toolbar
  background-color: rgba(251, 115, 4, 0.79)
</style>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useStore } from 'src/store';
import axios from 'axios';
import { api } from 'src/boot/axios';
import { Push } from 'src/lib/router/pushPage';

export default defineComponent({
  name: 'MainLayout',
  setup() {
    const store = useStore();
    const push = new Push();

    const loginInfoChecked = ref(false);
    const userID = ref(-1);
    const showName = ref('');
    const avatar = ref('');
    const grade = ref(-1);

    interface userDataJson {
      id: number;
      show_name: string;
      avatar: string;
      grade: number;
    }

    async function getUserInfo() {
      let path = '/me';
      await api
        .get(path)
        .then((res) => {
          const data = res.data as userDataJson;
          userID.value = data.id;
          showName.value = data.show_name;
          avatar.value = data.avatar;
          grade.value = data.grade;

          store.commit('user/setID', userID.value);
          store.commit('user/setGrade', grade.value);
          store.commit('user/setDataLoaded', true);
        })
        .catch((error) => {
          if (axios.isAxiosError(error)) {
            if (error.response) {
              switch (error.response.status) {
                // 400 -> 尚未登入
                case 400:
                case 404:
                  break;
              }
            }
          }
        });
      loginInfoChecked.value = true;
    }
    void getUserInfo();

    async function logout() {
      let path = '/logout';
      await api.post(path).then(() => {
        store.commit('user/setID', -1);
        store.commit('user/setGrade', -1);
        push.reload();
      });
    }

    return {
      loginInfoChecked,
      userID,
      showName,
      avatar,
      grade,
      push,
      logout,
    };
  },
});
</script>
