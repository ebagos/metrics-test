import { Client } from '@notionhq/client'
import { readFileSync } from 'fs'
import matter from 'gray-matter'
import { markdownToBlocks } from '@tryfabric/martian'

const core = require('@actions/core')
require('dotenv').config()
const token = process.env.NOTION_KEY
const databaseId = process.env.NOTION_DATABASE_ID
const filename = process.env.MARKDOWN_FILENAME
const tags = process.env.TAGS.split(',')
const title = process.env.PAGE_TITLE

console.log(databaseId, filename, tags, title)

async function main() {
  const notion = new Client({ auth: token })

  const notes = markdown_to_blocks(filename)
  
  try {
    await notion.pages.create({
      parent: { database_id: databaseId },
      properties: {
        Name: {
          type: 'title',
          title: [{ text: { content: title } }],
        },
        Tags: {
          type: 'multi_select',
          multi_select: tags.map(tag => ({ name: tag })),
        },
      },
      children: notes.body,
    })
  } catch (e) {
    console.error(`追加に失敗: `, e)
  }
}

function markdown_to_blocks(file) {
  const content = readFileSync(file)
  const matterResult = matter(content)

  return {
    body: markdownToBlocks(matterResult.content),
  }
}

main()
