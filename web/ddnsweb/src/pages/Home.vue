<template>
  <div>
    <div style="width: 20%; position: fixed; top: 5px; z-index: 9999; right: 10%">
      <el-autocomplete
        size="medium"
        style="float: right; width: 100%"
        v-model="searchKey"
        :fetch-suggestions="querySearch"
        placeholder="请输入搜索关键词"
        :trigger-on-focus="false"
        @select="inputHandleSelect"
        prefix-icon="el-icon-search"
      >
        <template slot-scope="{ item }">
          <div style="padding: 4px 0px 10px 0px">
            <el-breadcrumb v-if="item.type == 'WEB'" separator-class="el-icon-arrow-right">
              <el-breadcrumb-item style="color: #b0b0b0 !important; font-size: 16px">前端系统</el-breadcrumb-item>
              <el-breadcrumb-item style="color: #b0b0b0 !important; font-size: 16px" v-html="item.plat"></el-breadcrumb-item>
              <el-breadcrumb-item style="color: #b0b0b0 !important; font-size: 16px"></el-breadcrumb-item>
              <el-breadcrumb-item style="color: #b0b0b0 !important; font-size: 16px; margin-bottom: 10px" v-html="item.name"></el-breadcrumb-item>
            </el-breadcrumb>
            <el-breadcrumb v-if="item.type == 'API'">
              <el-breadcrumb-item style="color: #b0b0b0 !important; font-size: 16px" separator-class="el-icon-arrow-right">API</el-breadcrumb-item>
              <el-breadcrumb-item style="color: #b0b0b0 !important; font-size: 16px; margin-bottom: 10px" v-html="item.plat"></el-breadcrumb-item>
              <div
                v-for="(item1, index1) in item.server"
                :key="index1"
                style="margin-top: 5px; font-size: 16px; width: 100%; display: flex; align-items: center; flex-direction: row"
              >
                <div style="width: 95%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap">
                  {{ item1.systemName }}:&nbsp;&nbsp;&nbsp;<a target="_blank" v-bind:href="item1.address" class="systemurl" v-html="item1.address"> </a>
                </div>
                <i
                  class="el-icon-copy-document copy"
                  style="width: 5%"
                  v-clipboard:copy="item1.address"
                  v-clipboard:success="onCopy"
                  v-clipboard:error="onError"
                  icon="el-icon-search"
                ></i>
              </div>
            </el-breadcrumb>
          </div>
        </template>
      </el-autocomplete>
    </div>
    <el-container style="height: 95vh; border: 1px solid #eee; flex-direction: column; display: flex">
      <el-container>
        <el-aside width="20%" style="background-color: rgb(238, 241, 246)"> </el-aside>
        <el-main style="display: flex; align-items: center; background-color: white; margin-top: 15px; margin-bottom: 15px; flex-direction: column">
          <div style="width: 100%; padding-left: 20px; padding-right: 20px">
            <el-menu
              :default-active="activeIndex"
              class="el-menu-demo"
              mode="horizontal"
              background-color="#fff"
              active-text-color="#ff635f"
              text-color="#000"
              @select="handleSelect"
            >
              <el-menu-item index="web" style="width: 100px; font-size: 18px">前端系统</el-menu-item>
              <el-menu-item index="api" style="width: 100px; text-align: center; font-size: 18px">API</el-menu-item>
            </el-menu>
          </div>

          <div style="margin-top: 30px; width: 80%">
            <el-row :gutter="24" v-if="isWeb && !isAPI">
              <el-col
                :span="8"
                v-for="(item, index) in webTData"
                :key="index"
                class="el-col"
                style="align-items: center; display: flex; flex-direction: column; padding-left: 10px; height: 50px; margin-bottom: 20px; margin-top: 20px"
              >
                <div class="grid-content bg-purple keleyi">
                  <el-row :gutter="24" style="display: flex; align-items: center">
                    <el-col :span="8" class="el-col">
                      <el-image style="width: 50px; height: 50px" fit="contain" :src="item.imgurl">
                        <div slot="error">
                          <el-image :lazy="true" style="width: 50px; height: 50px" fit="cover" :src="item.imgdefaultrUrl"> </el-image>
                        </div>
                      </el-image>
                    </el-col>
                    <el-col :span="16" style="display: flex; align-items: flex-start; flex-direction: column">
                      <el-row :gutter="24">
                        <a target="_blank" v-bind:href="item.url" class="systemurl">
                          {{ item.name }}
                        </a>
                      </el-row>
                      <el-row :gutter="24" style="color: #999999">{{ item.plat }}</el-row>
                    </el-col>
                  </el-row>
                </div>
              </el-col>
            </el-row>
            <el-row :gutter="24" v-if="!isWeb && isAPI" style="margin-left: 5%; width: 100%">
              <el-col
                :span="8"
                v-for="(item, index) in apiTData"
                :key="index"
                class="el-col"
                style="align-items: center; display: flex; flex-direction: column; padding-left: 10px; height: 50px; margin-bottom: 20px; margin-top: 20px"
              >
                <div
                  class="grid-content bg-purple platOver"
                  @mouseover="apiSelectOver(index)"
                  @mouseleave="apiSelectLeave(index)"
                  style="display: flex; flex-direction: column; width: 100%"
                >
                  <el-row :gutter="24" style="display: flex; align-items: center; width: 100%">
                    <el-col :span="8" class="el-col">
                      <el-image style="width: 50px; height: 50px" fit="contain" :src="item.imgurl">
                        <div slot="error">
                          <el-image :lazy="true" style="width: 50px; height: 50px" fit="cover" :src="item.imgdefaultrUrl"> </el-image>
                        </div>
                      </el-image>
                    </el-col>
                    <el-col :span="16" style="display: flex; align-items: flex-start; flex-direction: column; padding-left: 0px">
                      <el-row :gutter="24" style="font-size: 16px; color: #111111">{{ item.plat }} </el-row>
                    </el-col>
                  </el-row>
                  <div v-if="isApiSelect == index" style="display: flex; align-items: center; flex-direction: column; width: 100%; padding: 0px 10px 10px 10px">
                    <div style="border-top: 1px solid #dddddd; width: 100%">
                      <div v-for="(item1, index1) in item.server" :key="index1" style="margin-top: 5px; width: 100%; display: flex; align-items: center; flex-direction: row">
                        <div style="width: 90%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #111111">
                          {{ item1.systemName }}:&nbsp;&nbsp;<a target="_blank" v-bind:href="item1.address" class="systemurl" style="font-size: 16px">
                            {{ item1.address }}
                          </a>
                        </div>
                        <i
                          class="el-icon-copy-document copy"
                          style="width: 10%"
                          v-clipboard:copy="item1.address"
                          v-clipboard:success="onCopy"
                          v-clipboard:error="onError"
                          icon="el-icon-search"
                        ></i>
                      </div>
                    </div>
                  </div>
                </div>
              </el-col>
            </el-row>
            <div style="display: flex; flex-direction: column; align-items: center">
              <div class="page-pagination">
                <el-pagination
                  background
                  @size-change="handleSizeChange"
                  @current-change="pageChange"
                  :current-page="paging.pi"
                  :page-size="paging.ps"
                  layout="total, prev, pager,next"
                  :total="totalCount"
                ></el-pagination>
              </div>
            </div>
          </div>
        </el-main>
        <el-aside width="20%" style="background-color: rgb(238, 241, 246)"> </el-aside>
      </el-container>
      <el-footer style="background-color: #e9e9e9; display: flex; flex-direction: row; align-items: center; text-align: center">
        <div style="width: 100%">
          <div><a>&nbsp;&nbsp;关于千行&nbsp;&nbsp;</a>|<a>&nbsp;&nbsp;联系我们&nbsp;&nbsp;</a>|<a>&nbsp;&nbsp;法律声明&nbsp;&nbsp;</a></div>
          <div style="color: #c0c0c0">{{ copyright }}</div>
        </div>
      </el-footer>
    </el-container>
  </div>
</template>

<script>
import "@/css/common.css";

let windowHeight = parseInt(window.innerHeight);
export default {
  name: "Home",
  watch: {
    search: {
      handler(val, oldVal) {
        this.searchVal(val)
      },
      immediate: false
    }
  },
  data() {
    return {
      isWeb: true,
      isAPI: false,
      isApiSelect: -1,
      search: "",
      activeIndex: 'web',
      webTData: [],
      apiTData: [],
      webData: [],
      apiData: [],
      copyright: "版权所有 @ 2014-" + new Date().getFullYear() + " 四川千行你我科技股份有限公司", //版权信息
      spanArr: [],//二维数组，用于存放单元格合并规则
      position: 0,//用于存储相同项的开始index
      paging: {
        ps: 21,
        pi: 1,
      },
      totalCount: 0,
      searchKey: "",
    };
  },
  created() {
    this.isWeb = true
    this.isAPI = false
    this.webData = []
    this.apiData = []
    this.queryData();
  },
  mounted() {
  },
  destroyed() {
  },
  methods: {
    pageChange: function (data) {
      this.paging.pi = data;
      this.queryData();
    },
    handleSizeChange(val) {
      this.paging.ps = val;
      this.queryData()
    },
    queryData() {
      this.$get("/ddns/query", {})
        .then(res => {
          var startIndex = (this.paging.pi - 1) * this.paging.ps
          var endIndex = (this.paging.pi) * this.paging.ps
          var imgIndex = 1
          if (res.web) {
            var web = res.web
            var webTempData = []
            var webIndex = 0
            for (var k = 0; k < web.length; k++) {
              var webClusters = web[k].clusters
              for (var key in webClusters) {
                var cluster = webClusters[key]
                var system = {
                  url: cluster.url,
                  imgurl: cluster.url + "/favicon.ico",
                  imgdefaultrUrl: require("../images/system" + imgIndex + ".png"),
                  name: cluster.system_cn_name,
                  plat: web[k].plat_cn_name,
                  type: "web",
                }
                webIndex++
                if (startIndex < webIndex && webIndex <= endIndex) {
                  webTempData.push(system)
                }
                imgIndex++
                if (imgIndex > 7) {
                  imgIndex = 1
                }
              }
            }
            this.webData = webTempData
            this.webTData = webTempData
            if (this.isWeb) {
              this.totalCount = webIndex
            }
          }
          if (res.api) {
            var api = res.api
            var apiTempData = []
            var apiIndex = 0
            var size = 0
            for (var k = 0; k < api.length; k++) {
              var apiClusters = api[k].clusters
              var tempPlat = {
                plat: api[k].plat_cn_name,
                //imgurl: cluster.url + "/favicon.ico",
                imgdefaultrUrl: require("../images/system" + imgIndex + ".png"),
                server: []
              }
              imgIndex++
              if (imgIndex > 7) {
                imgIndex = 1
              }
              for (var key in apiClusters) {
                var cluster = apiClusters[key]
                var server = {
                  url: cluster.url,
                  systemName: cluster.system_cn_name,
                  address: cluster.service_address,
                  type: "api",
                }
                tempPlat.server.push(server)
              }
              size++
              apiIndex++
              if (startIndex < apiIndex && apiIndex <= endIndex) {
                apiTempData.push(tempPlat)
              }
            }
            if (this.isAPI) {
              this.totalCount = size
            }
            this.apiTData = apiTempData
            this.apiData = apiTempData
          }
        })
        .catch(err => {
          console.log(err);
        });
    },
    handleSelect(key, keyPath) {
      if (key == "web") {
        this.isWeb = true
        this.isAPI = false
        this.queryData()
      }
      if (key == "api") {
        this.isWeb = false
        this.isAPI = true
        this.queryData()
      }
    },
    apiSelectOver(index) {
      this.isApiSelect = index
    },
    apiSelectLeave(index) {
      this.isApiSelect = -1
    },
    onCopy(e) { 　　 // 复制成功
      this.$message({
        message: '复制成功！',
        type: 'success'
      })
    },
    onError(e) {　　 // 复制失败
      this.$message({
        message: '复制失败！',
        type: 'error'
      })
    },
    searchData(searchVal) {
      var webData = this.webData.filter(item => ((~item.plat.indexOf(searchVal)) || (~item.name.indexOf(searchVal))));
      var apiData = []
      for (var i = 0; i < this.apiData.length; i++) {
        var apiItem = this.apiData[i]
        var temp = apiItem.server.filter(item => (~item.address.indexOf(searchVal)));
        if ((~apiItem.plat.indexOf(searchVal)) || temp.length > 0) {
          apiItem.server = temp
          apiItem.type = "API"
          apiData.push(apiItem)
        }
      }
      for (var i in webData) {
        webData[i].type = "WEB"
      }
      return JSON.parse(JSON.stringify(webData.concat(apiData)));
    },
    querySearch(searchVal, cb) {
      var tempResList = this.searchData(searchVal)
      let replaceReg = new RegExp(searchVal, "g");
      var replaceStr = '<span class="search-text">' + searchVal + "</span>";

      tempResList.forEach(element => {
        var platName = element.plat;
        element.plat = platName.replace(replaceReg, replaceStr);
        if (element.type == "WEB") {
          var name = element.name;
          element.name = name.replace(replaceReg, replaceStr);
        }
        if (element.type == "API") {
          element.server.forEach(item => {
            var address = item.address;
            item.address = address.replace(replaceReg, replaceStr);
          })
        }
      });
      // this.categoryDataDeal(tempResList);
      // 调用 callback 返回建议列表的数据
      cb(tempResList);
    },
    inputHandleSelect(item) {
      if (item.url) {
        window.location = item.url
      }
      if (item.address) {
        window.location = item.address
      }
    },
  },
  computed: {
  }
};
</script>


<style >
.search-text {
  color: #f56c6c;
}
.el-header {
  background-color: #b3c0d1;
  color: #333;
  line-height: 60px;
}
.el-menu-item.is-active {
  color: #fff;
}
.el-aside {
  color: #333;
}
.ivu-anchor-wrapper {
  z-index: 2;
  margin-left: 4px;
  height: 400px;
  max-height: 400px;
  width: 300px;
}

::-webkit-scrollbar {
  width: 0px;
}

.systemurl {
  color: #111111;
  font-size: 16px;
  transition: 0.5 s;
}
.systemurl:hover {
  color: #f56c6c;
}

.keleyi {
  width: 80%;
}
.keleyi:hover {
  background-color: #f8f8f8;
}
.platOver {
  width: 100%;
}
.platOver:hover {
  background-color: #f8f8f8;
  box-shadow: 1px 1px 3px #dddddd;
  z-index: 999;
}

.copy {
  color: #f56c6c;
}
.page-pagination {
  position: absolute;
  padding: 10px 15px;
  text-align: right;
  bottom: 80px;
}

/*当前页样式*/
.el-pagination.is-background .el-pager li:not(.disabled).active {
  background-color: #ffecec;
  color: #f56c6c;
  border: 1px solid #f56c6c;
}
/*当前页hover样式*/
.el-pagination.is-background .el-pager li:not(.disabled).active:hover {
  background-color: #ffecec;
  color: #f56c6c;
  border: 1px solid #f56c6c;
}
/*不是当前页其他页码hover样式*/
.el-pagination.is-background .el-pager li:not(.disabled):hover {
  color: #f56c6c;
}
</style>