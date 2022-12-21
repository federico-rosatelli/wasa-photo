<script>
export default {
    data:()=>{
        return {
            profiles: null,
            errormsg:null,
        }
    },
	props: ["imageUrl", "userNameF"],
}
</script>

<template>
    <div class="profile" v-if="imageUrl">
        <img v-bind:src="'http://localhost:3000'+imageUrl" alt="Profile Picture">
        <h4>{{userNameF}}</h4>
    </div>
</template>

<style>
.profile {
    display: flex;
}
.profile img{
    width: 5%;
    border-radius: 100%;
}
</style>
