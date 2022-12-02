<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/profile/199cd7ff-c6a1-437f-b50b-02df894ffbc5");
				this.some_data = response.data;
				console.log(this.some_data);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>

			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<div v-if="this.some_data" style="display: flex; gap: 40px;">
			<h2>{{this.some_data.Username}}</h2>
			<h2>{{this.some_data.Id}}</h2>
		</div>
		<div v-if="this.some_data" style="display: flex; gap: 20%;">
			<h2>{{this.some_data.Followers}}</h2>
			<h2>{{this.some_data.Followings}}</h2>
		</div>
		<div v-if="this.some_data" class="grid-container">
			<div v-for="item in this.some_data.Images" >
				<div class="grid-item">
					<img v-bind:src="'http://localhost:3000'+item.Location">
					<div style="display: flex;gap: 20%;">
						<h5>{{item.Text}}</h5>
						<h4>{{item.Comments}}</h4>
						<h4>{{item.Likes}}</h4>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.grid-container {
  display: grid;
  grid-template-columns: auto auto auto;
  padding: 10px;
}
.grid-item {
  padding: 20px;
  font-size: 30px;
  text-align: center;
  border: 1px solid #000000;
}
.grid-item img{
	object-fit: cover;
	width: 230px;
  	height: 230px;
}
</style>