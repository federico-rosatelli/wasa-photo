<script>
export default {
    data:()=>{
        return {
            profiles: null,
            errormsg:null,
        }
    },
	props: ["imageUrl"],
}
</script>

<template>
    <div class="profile" v-if="imageUrl">
        <img v-bind:src="'http://localhost:3000'+imageUrl" alt="Profile Picture">
    </div>
</template>

<style>
.profile img{
    width: 5%;
    height: 100%;
    border-radius: 100%;
}
</style>
