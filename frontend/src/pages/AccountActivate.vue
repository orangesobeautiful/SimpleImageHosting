<template>
  <q-page class="flex flex-center">
    <q-dialog v-model="showMsg">
      <q-card>
        <q-card-section>
          <div class="text-h6 text-center">{{ msgTitle }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none text-center">
          {{ msgContent }}
        </q-card-section>

        <q-card-actions align="center">
          <q-btn
            flat
            label="確認"
            color="primary"
            v-close-popup
            @click="confirmBtnClick"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
<script lang="ts">
import { ref, defineComponent } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { api } from 'boot/axios';
import { Push } from 'src/lib/router/pushPage';
export default defineComponent({
  name: 'AccountActivatePage',
  setup() {
    const route = useRoute();
    const push = new Push();

    const showMsg = ref(false);
    const actSuccess = ref(false);
    const msgTitle = ref('');
    const msgContent = ref('');
    const userID = ref(0);
    const grade = ref(0);
    const showName = ref('');

    interface userDataJson {
      user_id: number;
      grade: number;
      show_name: string;
    }

    async function activateAccount() {
      let token = route.params.token as string;
      let path = '/account-activate/' + token;
      await api
        .get(path)
        .then((res) => {
          let data = res.data as userDataJson;
          userID.value = data['user_id'];
          grade.value = data['grade'];
          showName.value = data['show_name'];
          actSuccess.value = true;
          msgTitle.value = '認證成功';
          msgContent.value = '歡迎您的加入 ' + showName.value;
          showMsg.value = true;
        })
        .catch((error) => {
          msgTitle.value = '認證失敗';
          if (axios.isAxiosError(error)) {
            if (error.response) {
              switch (error.response.status) {
                case 500:
                  msgContent.value = '伺服器內部錯誤';
                  showMsg.value = true;
                  break;
                case 403:
                  msgContent.value = '連結錯誤或已失效';
                  showMsg.value = true;
                  break;
              }
            }
          }
        });
    }
    void activateAccount();

    function confirmBtnClick() {
      if (actSuccess.value) {
        push.userImagesPage(userID.value);
      } else {
        push.homePage();
      }
    }

    return {
      showMsg,
      actSuccess,
      msgTitle,
      msgContent,
      userID,
      grade,
      showName,
      confirmBtnClick,
    };
  },
  methods: {},
});
</script>
