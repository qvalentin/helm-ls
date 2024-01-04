-- a minimal example config for setting up neovim with helm-ls and yamlls

-- setup lazy plugin manager
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
vim.g.mapleader = " "

require("lazy").setup({
  -- towolf/vim-helm provides basic syntax highlighting and filetype detection
  -- ft = 'helm' is important to not start yamlls
  --
  -- This is not needed when using treesitter (see below)
  -- { 'towolf/vim-helm',       ft = 'helm' },

  { "neovim/nvim-lspconfig", event = { "BufReadPre", "BufNewFile", "BufEnter" } },
  {
    "kmoschcau/tree-sitter-go-template",
    dependencies = "nvim-treesitter/nvim-treesitter",
    config = true,
  },
  {
    "nvim-treesitter/nvim-treesitter",
    build = ":TSUpdate",
    config = function()
      local configs = require("nvim-treesitter.configs")

      configs.setup({
        ensure_installed = { "gotmpl" },
        highlight = { enable = true },
        indent = { enable = true },
      })
    end

  }
})


local lspconfig = require('lspconfig')

-- setup helm-ls
lspconfig.helm_ls.setup {
  settings = {
    ['helm-ls'] = {
      yamlls = {
        path = "yaml-language-server",
      }
    }
  }
}

-- setup yamlls
lspconfig.yamlls.setup {}
