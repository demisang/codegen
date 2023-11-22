<template>
  <div class="container container-sm">
    <header class="d-flex flex-wrap py-3 mb-4">
      <ul class="nav nav-tabs">
        <li class="nav-item">
          <router-link exact-active-class="active" class="nav-link" to="/">Common templates</router-link>
        </li>
        <li class="nav-item">
          <router-link exact-active-class="active" class="nav-link" to="/about">About</router-link>
        </li>
      </ul>
    </header>

    <router-view/>

  </div>
</template>

<style lang="scss">

</style>
