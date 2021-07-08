<template>
<div >
    <div class="row justify-content-end pb-3 pr-3">
        <router-link :to="{name: 'add-output'}" class="btn btn-primary">Add Output</router-link>
    </div>
    <div class="row row-cols-1 row-cols-md-3">
        <OutputCard v-for="(output, index) in outputs"
            :key="index"
            :type="output.type"
            :name="output.name"
            :enable="output.enable">
        </OutputCard>
    </div>
</div>
</template>
<script>
import OutputCard from './OutputCard.vue';
import { mapState} from 'vuex'

export default {
    components: {
        OutputCard
    },
    computed: {
        ...mapState({
            outputs (state) {
                return state.outputs.all
            }
        })
    }

}
</script>