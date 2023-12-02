local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
  vim.fn.system({
    "git",
    "clone",
    "--filter=blob:none",
    "https://github.com/folke/lazy.nvim.git",
    "--branch=stable", -- latest stable release
    lazypath,
  })
end
vim.opt.rtp:prepend(lazypath)

vim.g.mapleader = " " -- Make sure to set `mapleader` before lazy so your mappings are correct

require("lazy").setup({
  { "sheerun/vim-polyglot",  lazy = false }, -- this contains 'towolf/vim-helm' but fixes some error, where yamlls would also be launched for helm files
  { "neovim/nvim-lspconfig", event = { "BufReadPre", "BufNewFile", "BufEnter" } }
})

local lspconfig = require('lspconfig')


lspconfig.helm_ls.setup {
  settings = {
    ['helm-ls'] = {
      yamlls = {
        path = "yaml-language-server",
      }
    }
  }
}

lspconfig.yamlls.setup {
  -- filetypes = vim.tbl_filter(function(ft)
  --   return not vim.tbl_contains({ "helm" }, ft)
  -- end, { 'yaml', 'yaml.docker-compose' }),
}
