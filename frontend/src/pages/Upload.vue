<template>
  <q-page class="flex flex-center">
    <loading-view :inputVisible="showLoading" />
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
            @click="push.homePage()"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
    <q-uploader
      v-if="showPage"
      class="upload-block"
      url="/api/image"
      accept="image/jpeg, image/png, image/gif"
      label="圖片上傳"
      method="POST"
      field-name="image[]"
      multiple
      bordered
    />
  </q-page>
</template>

<style lang="sass" scoped>
.upload-block
  width: 300px
  height: 600px
</style>

<script>
import { ref, defineComponent, watch } from 'vue';
import { useStore } from 'src/store';
import { Push } from 'src/lib/router/pushPage';
import LoadingView from '../components/LoadingView.vue';

export default defineComponent({
  name: 'PageUpload',
  components: {
    'loading-view': LoadingView,
  },
  setup() {
    const store = useStore();
    const userState = store.state.user;
    const push = new Push();

    const showLoading = ref(true);
    const showMsg = ref(false);
    const msgTitle = ref('');
    const msgContent = ref('');
    const showPage = ref(false);

    function pageLoad() {
      if (userState.dataLoaded) {
        showLoading.value = false;
        if (userState.id < 0) {
          // not login
          msgTitle.value = '您尚未登入';
          msgContent.value = '需要登入才能上傳圖片';
          showMsg.value = true;
        } else {
          showPage.value = true;
        }
      } else {
        showLoading.value = true;
      }
    }
    pageLoad();

    watch(userState.dataLoaded, () => {
      pageLoad();
    });

    return {
      showLoading,
      showMsg,
      msgTitle,
      msgContent,
      showPage,
      pageLoad,
      push,
    };
  },
});
</script>
