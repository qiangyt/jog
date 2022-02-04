<template>
  <q-layout view="hHh lpR fFf">
    <q-header elevated class="bg-primary text-white" height-hint="98">
      <q-toolbar class="GPLAY__toolbar">
        <q-btn dense flat round icon="menu" @click="toggleLeftDrawer" />

        <q-toolbar-title>
          <q-avatar>
            <img src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg" />
          </q-avatar>
          Jog
        </q-toolbar-title>

        <q-space />

        <div class="GPLAY__toolbar-input-container row no-wrap">
          <q-input
            dense
            outlined
            square
            v-model="search"
            placeholder="Search"
            class="bg-white col"
          />
          <q-btn
            class="GPLAY__toolbar-input-btn"
            color="primary"
            icon="search"
            unelevated
          />
        </div>

        <q-space />

        <q-btn dense flat round icon="menu" @click="toggleRightDrawer" />
      </q-toolbar>
    </q-header>

    <q-drawer
      behavior="desktop"
      show-if-above="false"
      v-model="leftDrawerOpen"
      side="left"
      bordered
    >
      <div><ul>
        <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/UserApiV1/authByEmailUsingPOST">User API #3.4 - 使用email登录认证</a></li>
        <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/listRingByOwnerUsingGET">Log API #1.5 - 列出一个拥有者（org）拥有的所有ring</a></li>
        <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/UserApiV1/getUserUsingGET">User API #2.1 - 指定id获取user</a></li>
      </ul></div>
    </q-drawer>

    <q-drawer
      behavior="mobile"
      show-if-above="false"
      elevated
      v-model="rightDrawerOpen"
      side="right"
      overlay
      bordered
    >
      <div>
        1. TODO: 获取当前user的所有未读消息<br />
      </div>
    </q-drawer>

    <q-page-container>
      <q-splitter v-model="splitterModel">
        <template v-slot:before>
          <!-- TODO：搬到单独的.vue里去 -->
          <q-timeline color="secondary">
            <q-timeline-entry title="timeline API">
              <div><ul>
                <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/listTimelineItemsUsingGET">Log API #2.2 - 列出timeline的item</a></li>
                <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/getTimelineItemUsingGET">Log API #3.1 - 获取timeline item</a></li>
                <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/createTimelineItemUsingPOST">Log API #2.3 - 在timeline里新建一个timeline item</a></li>
                <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/removeTimelineItemUsingDELETE">Log API #3.2 - 删除timeline item</a></li>
              </ul></div>
            </q-timeline-entry>

            <q-timeline-entry
              title="log entry #1"
              subtitle="2021-11-1 01:23:45"
            >
              <div>
                Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do
                eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
                enim ad minim veniam, quis nostrud exercitation ullamco laboris
                nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
                in reprehenderit in voluptate velit esse cillum dolore eu fugiat
                nulla pariatur. Excepteur sint occaecat cupidatat non proident,
                sunt in culpa qui officia deserunt mollit anim id est laborum.
              </div>
            </q-timeline-entry>

            <q-timeline-entry
              title="log entry #2"
              subtitle="2021-11-2 01:23:45"
            >
              <div>
                Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do
                eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
                enim ad minim veniam, quis nostrud exercitation ullamco laboris
                nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
                in reprehenderit in voluptate velit esse cillum dolore eu fugiat
                nulla pariatur. Excepteur sint occaecat cupidatat non proident,
                sunt in culpa qui officia deserunt mollit anim id est laborum.
              </div>
            </q-timeline-entry>

            <q-timeline-entry
              title="log entry #3"
              subtitle="2021-11-3 01:23:45"
            >
              <div>
                Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do
                eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
                enim ad minim veniam, quis nostrud exercitation ullamco laboris
                nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
                in reprehenderit in voluptate velit esse cillum dolore eu fugiat
                nulla pariatur. Excepteur sint occaecat cupidatat non proident,
                sunt in culpa qui officia deserunt mollit anim id est laborum.
              </div>
            </q-timeline-entry>
          </q-timeline>
        </template>

        <template v-slot:after>
          <q-tabs align="left">
            <q-route-tab to="/log1" label="log 1" />
            <q-route-tab to="/log2" label="log 2" />
            <q-route-tab to="/log3" label="log 3" />
          </q-tabs>
          <div><ul>
            <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/getLogFileUsingGET">Log API #4.1 - 指定id获取log file</a></li>
            <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/listLogItemsUsingGET">Log API #4.5 - 列出log file的item</a></li>
            <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/searchByKeywordUsingGET">Log API #6.1 - 使用简单关键字做全文检索</a></li>
            <li><a class="swagger_ui_link" href="http://i8700.wg:8080/swagger-ui/index.html#/LogApiV1/searchByExpressionUsingGET">Log API #6.2 - 使用检索表达式做全文检索</a></li>
          </ul></div>
          <router-view />
        </template>
      </q-splitter>
    </q-page-container>

    <q-footer elevated class="bg-grey-8 text-white">
      <q-toolbar>
        <q-toolbar-title>
          <div>Title</div>
        </q-toolbar-title>
      </q-toolbar>
    </q-footer>
  </q-layout>
</template>

<script>
import { ref } from 'vue';

export default {
  setup() {
    const leftDrawerOpen = ref(false);
    const rightDrawerOpen = ref(false);

    return {
      search: ref(''),
      leftDrawerOpen,
      toggleLeftDrawer() {
        leftDrawerOpen.value = !leftDrawerOpen.value;
      },

      rightDrawerOpen,
      toggleRightDrawer() {
        rightDrawerOpen.value = !rightDrawerOpen.value;
      },

      splitterModel: ref(38.2), // start at 38.2%
    };
  },
};
</script>

<style lang="sass">
 .swagger_ui_link
  color: green

.GPLAY
  &__toolbar
    height: 60px
  &__toolbar-input-container
    min-width: 100px
    width: 55%
  &__toolbar-input-btn
    border-radius: 0
    max-width: 60px
    width: 100%
</style>