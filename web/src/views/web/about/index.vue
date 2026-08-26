<template>

  <div class="about">


    <WebNavbar />


    <el-container class="main-content">


      <div class="container">


        <!-- 个人介绍 -->

        <section class="hero">


          <h1>
            FELIX
          </h1>


          <p>
            Full Stack Developer
          </p>


          <span>
            Technology · Creation · Exploration
          </span>


        </section>




        <!-- 关于我 -->


        <section class="card">


          <h2>
            ABOUT ME
          </h2>


          <p>
            我是一名专注于后端开发的程序员，
            喜欢研究 Web 技术、系统架构以及高性能服务。
          </p>


          <p>
            目前主要使用 Go 语言进行后端开发，
            同时掌握 Vue3 前端生态，
            致力于构建高性能、可扩展的 Web 应用。
          </p>


        </section>





        <!-- 技术栈 -->


        <section class="card">


          <h2>
            TECH STACK
          </h2>




          <div class="tech-grid">



            <div
                class="tech-item"
                v-for="item in techStack"
                :key="item.title"
            >


              <h3>
                {{item.title}}
              </h3>


              <p>
                {{item.desc}}
              </p>


              <div class="tags">


                <span
                    v-for="tag in item.tags"
                    :key="tag"
                >

                  {{tag}}

                </span>


              </div>



            </div>



          </div>


        </section>





        <!-- 项目能力 -->


        <section class="card">


          <h2>
            DEVELOPMENT EXPERIENCE
          </h2>


          <ul>


            <li>
              后端 API 服务开发
            </li>


            <li>
              RESTful API 设计
            </li>


            <li>
              MySQL 数据库设计与优化
            </li>


            <li>
              Redis 缓存设计
            </li>


            <li>
              Elasticsearch 搜索服务
            </li>


            <li>
              Linux + Docker 项目部署
            </li>


          </ul>


        </section>






      </div>


    </el-container>



  </div>


</template>





<script setup lang="ts">


import WebNavbar from "@/components/layout/WebNavbar.vue";



const techStack = [



  {


    title:"Backend",


    desc:"高性能后端服务开发",


    tags:[

      "Go",

      "Gin",

      "GORM",

      "JWT",

      "Redis",

      "MySQL",

      "Elasticsearch"

    ]

  },




  {


    title:"Frontend",


    desc:"现代化 Web 前端开发",


    tags:[

      "Vue3",

      "TypeScript",

      "Vite",

      "Pinia",

      "Element Plus",

      "ECharts"

    ]


  },





  {


    title:"DevOps",


    desc:"服务器部署与工程化",


    tags:[

      "Linux",

      "Docker",

      "Nginx",

      "Git",

      "CI/CD"

    ]


  }




];



</script>







<style scoped lang="scss">


.about{


  min-height:100vh;


  background:#050505;



  color:white;




}




.main-content{


  margin-top:70px;


  display:flex;


  justify-content:center;



}





.container{


  width:100%;


  max-width:1200px;


  padding:40px;



}







.hero{


  text-align:center;


  padding:80px 0;


  animation:fade .8s ease;



  h1{


    font-size:80px;


    letter-spacing:20px;


    margin:0;



    background:

    linear-gradient(
        90deg,
        white,
        #666
    );


    -webkit-background-clip:text;


    color:transparent;



  }



  p{


    margin-top:20px;


    font-size:28px;


    color:#aaa;



  }



  span{


    color:#666;


    letter-spacing:5px;


  }



}








.card{


  margin-bottom:40px;


  padding:40px;



  border-radius:20px;



  background:

  rgba(255,255,255,.04);



  border:

  1px solid rgba(255,255,255,.1);



  backdrop-filter:

  blur(20px);



  animation:

  fadeUp .8s ease;



  transition:.3s;




  &:hover{


    transform:

    translateY(-8px);



    border-color:

    rgba(255,255,255,.4);



  }



  h2{


    font-size:28px;


    letter-spacing:6px;


    margin-bottom:30px;


  }



  p{


    color:#aaa;


    line-height:2;


    font-size:16px;


  }



}







.tech-grid{


  display:grid;


  grid-template-columns:

  repeat(3,1fr);


  gap:25px;


}







.tech-item{


  padding:30px;


  background:#0d0d0d;


  border-radius:15px;



  transition:.3s;



  &:hover{


    transform:

    translateY(-10px);



  }



  h3{


    font-size:22px;


  }




}







.tags{


  display:flex;


  flex-wrap:wrap;


  gap:10px;



  span{


    padding:

    6px 14px;


    border-radius:20px;


    background:#181818;


    border:

    1px solid #333;


    color:#ccc;


    font-size:13px;



  }



}






ul{


  padding-left:20px;


  color:#aaa;


  line-height:2.5;



}








@keyframes fade{


  from{

    opacity:0;

  }


  to{

    opacity:1;

  }

}




@keyframes fadeUp{


  from{


    opacity:0;


    transform:

    translateY(40px);


  }


  to{


    opacity:1;


    transform:

    translateY(0);



  }


}






@media(max-width:900px){


  .tech-grid{


    grid-template-columns:

    1fr;


  }



  .hero h1{


    font-size:50px;


  }


}




</style>