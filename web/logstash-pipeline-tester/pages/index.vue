<template>
  <v-row justify="center" align="center" class="mt-2">
    <v-col cols="12" sm="10">
      <v-card elevation="2">
        <v-card-title>Please fill the form</v-card-title>
        <v-card-text>Your request uuid: <kbd>{{ uuid }}</kbd></v-card-text>
        <v-divider class="mx-4"></v-divider>
        <v-form v-model="valid" @submit.prevent>
          <v-container>
            <div class="text-h6">Logstash Pipeline</div>
            <div class="text-caption mb-4">Please fill one of using the textbox or upload your pipeline file.</div>
            <!-- <v-row dense>
              <v-col>
                <v-file-input outlined truncate-length="15" label="Logstash pipeline file"></v-file-input>
              </v-col>
            </v-row> -->
            <v-row dense>
              <v-col>
                <v-textarea outlined name="pipelineInput" v-model="pipelineInput" :rules="pipelineInputRules"
                  label="Logstash Pipeline"></v-textarea>
              </v-col>
            </v-row>

            <div class="text-h6">Sample logs</div>
            <div class="text-caption mb-4">Please fill one of using the textbox or upload your log file.</div>
            <!-- <v-row dense>
              <v-col>
                <v-file-input outlined truncate-length="15" label="Logs input file"></v-file-input>
              </v-col>
            </v-row> -->
            <v-row dense>
              <v-col>
                <v-textarea outlined name="logInput" v-model="logInput" :rules="logInputRules" label="Logs input">
                </v-textarea>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-btn color="info" class="float-right" @click="submit" :disabled="!valid || loading">
                  submit
                </v-btn>
                <v-btn @click="clear" class="mr-4 float-right" :disabled="loading">
                  clear
                </v-btn>
              </v-col>
            </v-row>
          </v-container>
        </v-form>
        <v-progress-linear v-if="loading" indeterminate color="red"></v-progress-linear>
      </v-card>
    </v-col>
    <v-col cols="12" sm="10" v-show="resultOutput">
      <v-card>
        <v-card-title>Result</v-card-title>
        <v-card-text>message: {{ resultMessage }}</v-card-text>
        <v-list>
          <v-list-group color="red">
            <template v-slot:activator>
              <v-list-item-title>Logstash Pipeline</v-list-item-title>
            </template>
            <v-list-item>
              <v-list-item-content>
                <v-card v-bind:class="[$vuetify.theme.dark ? 'grey lighten-4' : 'grey darken-4' ]">
                  <v-card-text>
                    <pre v-bind:class="[$vuetify.theme.dark ? 'grey--text text--darken-4' : 'white--text' ]">
{{ this.pipelineInput }}
            </pre>
                  </v-card-text>
                </v-card>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
        </v-list>
        <v-list>
          <v-list-group color="red">
            <template v-slot:activator>
              <v-list-item-title>Sample Logs</v-list-item-title>
            </template>
            <v-list-item>
              <v-list-item-content>
                <v-card v-bind:class="[$vuetify.theme.dark ? 'grey lighten-4' : 'grey darken-4' ]">
                  <v-card-text>
                    <pre v-bind:class="[$vuetify.theme.dark ? 'grey--text text--darken-4' : 'white--text' ]">
{{ this.logInput }}
            </pre>
                  </v-card-text>
                </v-card>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
        </v-list>
        <v-list>
          <v-list-group color="red" :value="true">
            <template v-slot:activator>
              <v-list-item-title>Result Output</v-list-item-title>
            </template>
            <v-list-item>
              <v-list-item-content>
                <v-card v-bind:class="[$vuetify.theme.dark ? 'grey lighten-4' : 'grey darken-4' ]">
                  <v-card-text>
                    <pre v-bind:class="[$vuetify.theme.dark ? 'grey--text text--darken-4' : 'white--text' ]">
{{ this.resultOutput }}
            </pre>
                  </v-card-text>
                </v-card>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
        </v-list>
      </v-card>
    </v-col>
  </v-row>
</template>
<style>
  pre {
    /* white-space: pre-wrap; */
    overflow: auto;
  }

</style>
<script>
  var qs = require('qs');

  export default {
    data() {
      return {
        loading: false,
        valid: false,
        uuid: "",
        pipelineInput: "",
        pipelineInputRules: [
          v => !!v || 'Logstash pipeline is required'
        ],
        logInput: "",
        logInputRules: [
          v => !!v || 'Log Sample is required'
        ],
        resultOutput: "",
        resultMessage: ""
      }
    },
    mounted: function () {
      this.uuid = this.uuidv4()
    },
    methods: {
      uuidv4() {
        return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
          (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
        );
      },
      submit() {
        // reset the output and message
        this.resultOutput = ""
        this.resultMessage = ""

        let data = {
          "uuid": this.uuid,
          "pipeline_input": this.pipelineInput,
          "log_input": this.logInput
        }
        this.submitPipeline(data)
      },
      clear() {
        this.pipelineInput = ""
        this.logInput = ""
        this.resultOutput = ""
        this.resultMessage = ""
      },
      async submitPipeline(req) {
        this.loading = true
        const data = await this.$axios.$post('http://localhost:8080/api/pipeline', qs.stringify(req), {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          }
        })
        this.resultOutput = data.data.output
        this.resultMessage = data.message
        this.loading = false
      }
    }
  }

</script>
