return {
	{
		'nvimtools/none-ls.nvim',
		dependencies = { "nvim-lua/plenary.nvim" },
		config = function()
			local null_ls    = require("null-ls")
			local helpers    = require("null-ls.helpers")
			local methods    = require("null-ls.methods")
			local FORMATTING = methods.internal.FORMATTING

			local nasm_fmt   = helpers.make_builtin({
				name = "nasmfmt",
				meta = { description = "NASM formatter" },
				method = FORMATTING,
				filetypes = { "asm", "nasm" },

				generator_opts = {
					command = "nasm-fmt",
					args = { "$FILENAME" },
					to_stdin = false,
					to_temp_file = true,
				},

				factory = helpers.formatter_factory,
			})

			null_ls.setup({
				sources = { nasm_fmt },
			})
		end
	},
}
