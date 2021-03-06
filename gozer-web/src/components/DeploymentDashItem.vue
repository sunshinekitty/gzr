<template>
  <div class="panel panel-default application-card">
    <div class="panel-heading">
      <h3 class="panel-title">
        <router-link :to="'/deployment/' + deploymentInternal.name">{{deploymentInternal.name}}
        </router-link>
        <span v-if="deploymentImageAge" class="label label-default pull-right">{{deploymentImageAge}}</span>
      </h3>
      <sub v-if="!!deploymentAppImage">{{deploymentAppImage.name}}</sub>
    </div>
    <div class="panel-body">
      <div v-if="!!deploymentAppImage">
        <a class="btn btn-default btn-github"
           :href="deploymentAppImage.metadata['git-hub-base'] + '/commit/' + deploymentAppImage.metadata['git-commit']"><img
          src="../img/GitHub-Mark-32px.png" alt="github"/></a>
        <span class="label label-info label-github" v-for="tag in deploymentAppImage.metadata['git-tag']">
          <a :href="deploymentAppImage.metadata['git-hub-base'] + '/tree/' + tag">{{tag}}</a>
        </span>
      </div>
      <div v-else-if="loading">
        <spinner v-model="loading" size="lg"></spinner>
      </div>
      <div v-else class="text-center middle">
        <p>No metadata for deployment</p>
      </div>
    </div>
  </div>
</template>

<script>
  import {accordion, panel, spinner} from 'vue-strap'
  import moment from 'moment';
  import Promise from 'bluebird';
  import imagesService from '../services/ImagesService'

  function cleanUpGitLink(gitLink){
      return gitLink
        .replace('git@github.com:', 'https://github.com/')
        .replace('.git', '')
  }

  export default {
    props: {deployment: {type: Object}},
    data() {
      return {
        deploymentImageAge: undefined,
        deploymentAppImage: undefined,
        deploymentInternal: JSON.parse(JSON.stringify(this.deployment)),
      };
    },
    created: function () {
      const vm = this;
      vm.loading = true;
      Promise.map(vm.deploymentInternal.containers, function (container) {
        let containerImageParts = container.image.split(':');
        let containerName = containerImageParts[0];
        let containerVersion = containerImageParts[1];
        return imagesService.getByVersion(containerName, containerVersion)
          .then(function (image) {
            if (image) {
              image['git-hub-base'] = cleanUpGitLink(image['git-origin'])
              vm.deploymentAppImage = image;

              vm.deploymentImageAge = moment(image.metadata["created-at"]).fromNow();

            }
          })
          .catch(function (error) {
            if (error && error.status !== 404) {
              console.log(error)
            }
          })
      })
        .finally(function () {
          vm.loading = false;
        })
      ;
    },
    components: {
      accordion,
      panel,
      spinner
    }
  };
</script>

<style scoped>

  .btn-github {
    padding: 2px 6px;
    margin: 5px;
  }

  .btn-github img {
    height: 28px;
    width: 28px;
  }

  .label-github a {
    color: white;
  }

  .application-card {
    height: 100%;
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 133px;
    overflow: hidden;
    text-overflow: ellipsis;
    flex-basis: 320px;
    min-width: 320px;
  }

  .panel-heading {
    flex: 1;
  }

  .panel-body {
    flex: 2;
  }

</style>
